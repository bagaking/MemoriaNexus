package model

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/bagaking/goulp/wlog"
	"github.com/khgame/memstore/cachekey"
	"github.com/khicago/got/util/delegate"
	"github.com/khicago/got/util/typer"
	"github.com/khicago/irr"

	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/internal/utils/cache"
)

const (
	MaxRetryAttempts      = 3
	TagCacheExpiration    = 24 * time.Hour
	CacheInvalidationFlag = "invalid"
)

type (
	EntityType uint8

	Tag struct {
		UserID     utils.UInt64 `gorm:"primaryKey"`
		Tag        string       `gorm:"primaryKey"`
		EntityID   utils.UInt64 `gorm:"primaryKey"`
		EntityType EntityType   `gorm:"not null"`
		CreatedAt  time.Time
		UpdatedAt  time.Time
		DeletedAt  gorm.DeletedAt
	}

	ParamUserTag struct {
		UserID utils.UInt64
		Tag    string
	}

	ParamUserTagType struct {
		UserID utils.UInt64 `cachekey:"user_id"`
		Tag    string       `cachekey:"tag"`
		Type   EntityType   `cachekey:"entity_type"`
	}
)

var (
	CKEntity2Tags          = cachekey.MustNewSchema[utils.UInt64]("entity:{entity_id}:tags", TagCacheExpiration)
	CKUser2Tags            = cachekey.MustNewSchema[utils.UInt64]("entity:{user_id}:tags", TagCacheExpiration)
	CKUserTag2Entities     = cachekey.MustNewSchema[ParamUserTag]("user:{user_id}:tag:{tag}:entities", TagCacheExpiration)
	CKUserTagType2Entities = cachekey.MustNewSchema[ParamUserTagType]("user:{user_id}:tag:{tag}:type:{entity_type}:entities", TagCacheExpiration)
	CKTag2Users            = cachekey.MustNewSchema[string]("tag:{tag}:users", TagCacheExpiration)
)

const (
	EntityTypeItem    EntityType = 1
	EntityTypeBook    EntityType = 2
	EntityTypeDungeon EntityType = 3
)

func (Tag) TableName() string {
	return "tags"
}

// GetTagsByEntity retrieves tags associated with a given entity ID.
func GetTagsByEntity(ctx context.Context, tx *gorm.DB, entityID utils.UInt64) ([]string, error) {
	log := wlog.ByCtx(ctx, "GetTagsByEntity")
	cacheKey := CKEntity2Tags.MustBuild(entityID)
	cachedTags, err := cache.SET().GetAll(ctx, cacheKey)
	if err == nil && len(cachedTags) > 0 && cachedTags[0] != CacheInvalidationFlag {
		log.Debugf("got cached tags: %v", cachedTags)
		return cachedTags, nil
	}

	var tags []string
	if err = tx.Model(&Tag{}).Where("entity_id = ? AND deleted_at IS NULL", entityID).Pluck("tag", &tags).Error; err != nil {
		return nil, irr.Wrap(err, "failed to get tags by entity")
	}
	log.Debugf("got tags from db: %v", tags)

	if len(tags) > 0 {
		if err = cache.SET().Insert(ctx, cacheKey, TagCacheExpiration, tags...); err != nil {
			log.Warnf("failed to insert tags into cache, key= %v, tags= %v", cacheKey, tags)
		}
	}
	log.Trace("inserted tags into cache success, key= %v, tags= %v", cacheKey, tags)

	return tags, nil
}

// AddEntityTags adds tags to an entity for a user.
func AddEntityTags(ctx context.Context, tx *gorm.DB, userID utils.UInt64, entityType EntityType, entityID utils.UInt64, tags ...string) error {
	uts := make([]Tag, 0, len(tags))
	for _, tag := range tags {
		userTag := Tag{
			UserID:     userID,
			Tag:        tag,
			EntityID:   entityID,
			EntityType: entityType,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		uts = append(uts, userTag)
	}

	if err := tx.Create(&uts).Error; err != nil {
		return irr.Wrap(err, "failed to add entity tag")
	}

	// Invalidate relevant caches
	if err := invalidateCacheForEntity(ctx, entityID); err != nil {
		return err
	}

	// all user's tag-related cache should be cleared
	// todo: consider to clear only the cache with the entity type
	if err := invalidateCacheForTag(ctx, tx, tags...); err != nil {
		return err
	}
	return nil
}

// RemoveEntityTags removes tags from an entity for a user.
func RemoveEntityTags(ctx context.Context, tx *gorm.DB, userID utils.UInt64, entityID utils.UInt64, tagsToRemove ...string) error {
	if err := tx.Model(&Tag{}).Where("user_id = ? AND tag IN ? AND entity_id = ?", userID, tagsToRemove, entityID).Update("deleted_at", time.Now()).Error; err != nil {
		return irr.Wrap(err, "failed to remove entity tag")
	}

	// Invalidate relevant caches
	if err := invalidateCacheForEntity(ctx, entityID); err != nil {
		return err
	}
	if err := invalidateCacheForTag(ctx, tx, tagsToRemove...); err != nil {
		return err
	}

	return nil
}

// UpdateEntityTagsDiff handles the difference between the tags of an entity and the new tags.
func UpdateEntityTagsDiff(ctx context.Context, tx *gorm.DB, userID utils.UInt64, entityID utils.UInt64, tags []string) error {
	tagsExist, err := GetTagsByEntity(ctx, tx, entityID)
	if err != nil {
		return irr.Wrap(err, "failed to get exist tags")
	}

	toAdd, toRemove := diffLists(tagsExist, tags)
	if err = AddEntityTags(ctx, tx, userID, EntityTypeItem, entityID, toAdd...); err != nil {
		return irr.Wrap(err, "failed to add tags")
	}
	if err = RemoveEntityTags(ctx, tx, userID, entityID, toRemove...); err != nil {
		return irr.Wrap(err, "failed to remove tags")
	}
	return nil
}

// GetEntities retrieves entity IDs associated with a given user ID and tag.
func GetEntities(ctx context.Context, tx *gorm.DB, userID utils.UInt64, tag string) ([]utils.UInt64, error) {
	cacheKey := CKUserTag2Entities.MustBuild(ParamUserTag{UserID: userID, Tag: tag})
	cachedEntities, err := cache.SET().GetAllUInt64s(ctx, cacheKey)
	if err == nil && len(cachedEntities) > 0 && cachedEntities[0] != utils.UInt64(0) {
		return cachedEntities, nil
	}

	var entityIDs []utils.UInt64
	if err = tx.Model(&Tag{}).Where("user_id = ? AND tag = ? AND deleted_at IS NULL", userID, tag).Pluck("entity_id", &entityIDs).Error; err != nil {
		return nil, irr.Wrap(err, "failed to get entities by user and tag")
	}

	if err = cache.SET().InsertUInt64s(ctx, cacheKey, TagCacheExpiration, entityIDs...); err != nil {
		return nil, err
	}

	return entityIDs, nil
}

// GetEntitiesOfType retrieves entity IDs associated with a given user ID, tag, and entity type.
func GetEntitiesOfType(ctx context.Context, tx *gorm.DB, userID utils.UInt64, tag string, entityType EntityType) ([]utils.UInt64, error) {
	cacheKey := CKUserTagType2Entities.MustBuild(ParamUserTagType{UserID: userID, Tag: tag, Type: entityType})
	cachedEntities, err := cache.SET().GetAllUInt64s(ctx, cacheKey)
	if err == nil && len(cachedEntities) > 0 && cachedEntities[0] != utils.UInt64(0) {
		return cachedEntities, nil
	}

	var entityIDs []utils.UInt64
	if err = tx.Model(&Tag{}).Where("user_id = ? AND tag = ? AND entity_type = ? AND deleted_at IS NULL", userID, tag, entityType).Pluck("entity_id", &entityIDs).Error; err != nil {
		return nil, irr.Wrap(err, "failed to get entities by user, tag, and entity type")
	}

	if err = cache.SET().InsertUInt64s(ctx, cacheKey, TagCacheExpiration, entityIDs...); err != nil {
		return nil, err
	}

	return entityIDs, nil
}

// GetTagsByUser retrieves tags associated with a given user ID.
func GetTagsByUser(ctx context.Context, tx *gorm.DB, userID utils.UInt64) ([]string, error) {
	cacheKey := CKUser2Tags.MustBuild(userID)
	cachedTags, err := cache.SET().GetAll(ctx, cacheKey)
	if err == nil && len(cachedTags) > 0 && cachedTags[0] != CacheInvalidationFlag {
		return cachedTags, nil
	}

	var tags []string
	if err := tx.Model(&Tag{}).Where("user_id = ? AND deleted_at IS NULL").Distinct().Pluck("tag", &tags).Error; err != nil {
		return nil, irr.Wrap(err, "failed to get tags by user")
	}

	if err := cache.SET().Insert(ctx, cacheKey, TagCacheExpiration, tags...); err != nil {
		return nil, err
	}

	return tags, nil
}

// GetUsersByTag retrieves user IDs associated with a given tag.
// if the tx is nil, it will only use cache.
func GetUsersByTag(ctx context.Context, tx *gorm.DB, tag string) ([]utils.UInt64, error) {
	log := wlog.ByCtx(ctx, "GetUsersByTag").WithField("tag", tag)
	cacheKey := CKTag2Users.MustBuild(tag)
	cachedUsers, err := cache.SET().GetAllUInt64s(ctx, cacheKey)
	log.Debugf("got cached users: %v", cachedUsers)
	if err == nil && len(cachedUsers) > 0 && cachedUsers[0] != utils.UInt64(0) {
		return cachedUsers, nil
	}

	var userIDs []utils.UInt64
	if err = tx.Model(&Tag{}).
		Where("tag = ? AND deleted_at IS NULL", tag).
		Distinct().
		Pluck("user_id", &userIDs).
		Error; err != nil {
		return nil, irr.Wrap(err, "failed to get users by tag %s", tag)
	}
	log.Debugf("got users from db: %v", userIDs)

	if err = cache.SET().InsertUInt64s(ctx, cacheKey, TagCacheExpiration, userIDs...); err != nil {
		return nil, err
	}
	log.Debugf("inserted users into cache success, key= %v, users= %v", cacheKey, userIDs)

	return userIDs, nil
}

// RenameTag updates the tag name and clears relevant caches.
func RenameTag(ctx context.Context, tx *gorm.DB, userID utils.UInt64, oldTag, newTag string) error {
	// 查询与旧标签关联的所有实体
	entityIDs, err := GetEntities(ctx, tx, userID, oldTag)
	if err != nil {
		return irr.Wrap(err, "failed to get entities by old tag")
	}

	// 检查这些实体ID是否已经与新标签关联
	var existingTags []Tag
	if err := tx.Where("user_id = ? AND tag = ? AND entity_id IN ?", userID, newTag, entityIDs).Find(&existingTags).Error; err != nil {
		return irr.Wrap(err, "failed to find existing tags")
	}

	existingTagsMap := make(map[utils.UInt64]Tag)
	for _, tag := range existingTags {
		existingTagsMap[tag.EntityID] = tag
	}

	var tagsToCreate []Tag
	var tagsToUpdate []Tag
	for _, entityID := range entityIDs {
		if _, exists := existingTagsMap[entityID]; !exists {
			tagsToCreate = append(tagsToCreate, Tag{
				UserID:     userID,
				Tag:        newTag,
				EntityID:   entityID,
				EntityType: EntityTypeItem, // 根据实际情况调整
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			})
		} else if tag, ok := existingTagsMap[entityID]; ok && tag.DeletedAt.Valid {
			tagsToUpdate = append(tagsToUpdate, tag)
		}
	}

	// 批量创建新的标签关联
	if len(tagsToCreate) > 0 {
		if err := tx.Create(&tagsToCreate).Error; err != nil {
			return irr.Wrap(err, "failed to batch create new tag associations")
		}
	}

	// 批量更新已软删除的标签关联
	if err = tx.Model(&Tag{}).Where(
		"user_id = ? AND tag = ? AND entity_id IN ?",
		userID, newTag, typer.SliceMap(tagsToUpdate, func(from Tag) utils.UInt64 {
			return from.EntityID
		})).
		Updates(Tag{DeletedAt: gorm.DeletedAt{}}).Error; err != nil {
		return irr.Wrap(err, "failed to batch update soft-deleted tag associations")
	}

	// 异步清除与旧标签相关的缓存
	go func() {
		if err := deleteOldTagData(ctx, tx, oldTag); err != nil {
			// todo: 记录日志或采取其他措施
		}
	}()

	return nil
}

// deleteOldTagData deletes old tag data and clears relevant caches.
func deleteOldTagData(ctx context.Context, tx *gorm.DB, oldTag string) error {
	// 将旧记录软删
	if err := tx.Model(&Tag{}).Where("tag = ?", oldTag).Update("deleted_at", time.Now()).Error; err != nil {
		return irr.Wrap(err, "failed to soft delete old tag")
	}

	// 清除缓存
	if err := invalidateCacheForTag(ctx, tx, oldTag); err != nil {
		return irr.Wrap(err, "failed to invalidate cache for old tag")
	}

	// 清除与 entity 相关的缓存
	entityIDs, err := GetEntities(ctx, tx, 0, oldTag)
	if err != nil {
		return irr.Wrap(err, "failed to get entities by old tag")
	}
	for _, entityID := range entityIDs {
		if err := invalidateCacheForEntity(ctx, entityID); err != nil {
			return irr.Wrap(err, "failed to invalidate cache for entity")
		}
	}

	return nil
}

// invalidateCacheForTag invalidates caches related to a given tag.
// if the tx is nil, it will only use cache.
// but for clear all user's cache, its better to use tx not nil, to make sure the cache is cleared.
func invalidateCacheForTag(ctx context.Context, tx *gorm.DB, tags ...string) error {
	for _, tag := range tags {
		// find all users with the tag
		userIDs, err := GetUsersByTag(ctx, tx, tag)
		if err != nil {
			return err
		}

		// clear caches for each user
		for _, userID := range userIDs {
			if err = cache.SET().Clear(ctx, CKUser2Tags.MustBuild(userID),
				MaxRetryAttempts); err != nil {
				return err
			}
			if err = cache.SET().Clear(ctx, CKUserTag2Entities.MustBuild(ParamUserTag{UserID: userID, Tag: tag}),
				MaxRetryAttempts); err != nil {
				return err
			}
			if err = cache.SET().Clear(ctx, CKUserTagType2Entities.MustBuild(ParamUserTagType{UserID: userID, Tag: tag, Type: EntityTypeItem}),
				MaxRetryAttempts); err != nil {
				return err
			}
			if err = cache.SET().Clear(ctx, CKUserTagType2Entities.MustBuild(ParamUserTagType{UserID: userID, Tag: tag, Type: EntityTypeBook}),
				MaxRetryAttempts); err != nil {
				return err
			}
			if err = cache.SET().Clear(ctx, CKUserTagType2Entities.MustBuild(ParamUserTagType{UserID: userID, Tag: tag, Type: EntityTypeDungeon}),
				MaxRetryAttempts); err != nil {
				return err
			}
		}

		// clear tag_to_users cache after all user's cache is cleared，cuz the cache is used to find all users with the tag
		if err = cache.SET().Clear(ctx, CKTag2Users.MustBuild(tag), MaxRetryAttempts); err != nil {
			return err
		}
	}
	return nil
}

// invalidateCacheForEntity invalidates caches related to a given entity ID.
func invalidateCacheForEntity(ctx context.Context, entityID utils.UInt64) error {
	// there are only entity tags cache to clear
	return cache.SET().Clear(ctx, CKEntity2Tags.MustBuild(entityID), MaxRetryAttempts)
}

// diffLists takes two slices of strings and returns two slices:
// one with elements to add (present in newList but not in oldList)
// and one with elements to remove (present in oldList but not in newList).
func diffLists(oldList, newList []string) (toAdd, toRemove []string) {
	// Threshold to decide between loop and map approach
	threshold := 10
	var nExistInOld delegate.Predicate[string]
	var nExistInNew delegate.Predicate[string]

	type TrueTable = map[string]struct{}
	// Use map-based approach for larger lists for better performance
	if len(oldList) > threshold {
		trueTable := typer.SliceReduce(oldList, func(v string, target TrueTable) TrueTable {
			target[v] = struct{}{}
			return target
		}, make(TrueTable))
		nExistInOld = func(v string) bool {
			_, ok := trueTable[v]
			return !ok
		}
	} else {
		nExistInOld = func(v string) bool {
			return typer.SliceContains(oldList, v)
		}
	}

	if len(newList) > threshold {
		trueTable := typer.SliceReduce(newList, func(v string, target TrueTable) TrueTable {
			target[v] = struct{}{}
			return target
		}, make(TrueTable))
		nExistInNew = func(v string) bool {
			_, ok := trueTable[v]
			return !ok
		}
	} else {
		nExistInOld = func(v string) bool {
			return typer.SliceContains(newList, v)
		}
	}

	toAdd = typer.SliceFilter(newList, nExistInOld)
	toRemove = typer.SliceFilter(oldList, nExistInNew)
	return toAdd, toRemove
}

// FindItemsOfTag returns the items
func FindItemsOfTag(ctx context.Context, tx *gorm.DB, userID utils.UInt64, tag string, pager *utils.Pager) ([]Item, error) {
	itemIDs, err := GetEntitiesOfType(ctx, tx, userID, tag, EntityTypeItem)
	if err != nil {
		return nil, irr.Wrap(err, "failed to find items of tag= %v", tag)
	}

	var items []Item
	if err = tx.Where("id IN ?", itemIDs).Offset(pager.Offset).Limit(pager.Limit).Find(&items).Error; err != nil {
		return nil, irr.Wrap(err, "failed to fetch items by id in %v", itemIDs)
	}
	return items, nil
}

// FindBooksOfTag returns the items
func FindBooksOfTag(ctx context.Context, tx *gorm.DB, userID utils.UInt64, tag string, pager *utils.Pager) ([]Book, error) {
	bookIDs, err := GetEntitiesOfType(ctx, tx, userID, tag, EntityTypeItem)
	if err != nil {
		return nil, irr.Wrap(err, "failed to find items of tag= %v", tag)
	}
	var books []Book
	if err = tx.Where("id IN ?", bookIDs).Offset(pager.Offset).Limit(pager.Limit).Find(&books).Error; err != nil {
		return nil, irr.Wrap(err, "failed to fetch books by id in %v", bookIDs)
	}
	return books, nil
}
