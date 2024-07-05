package model

import (
	"context"
	"time"

	"github.com/khicago/got/util/delegate"
	"github.com/khicago/got/util/typer"
	"github.com/khicago/irr"
	"gorm.io/gorm"

	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/internal/utils/cache"
)

const (
	MaxRetryAttempts      = 3
	CacheExpiration       = 24 * time.Hour
	CacheInvalidationFlag = "invalid"
)

type (
	EntityType uint8

	UserTag struct {
		UserID     utils.UInt64 `gorm:"primaryKey"`
		Tag        string       `gorm:"primaryKey"`
		EntityID   utils.UInt64 `gorm:"primaryKey"`
		EntityType EntityType   `gorm:"not null"`
		CreatedAt  time.Time
		UpdatedAt  time.Time
		DeletedAt  gorm.DeletedAt
	}

	ParamUserTag struct {
		UserID utils.UInt64 `json:"user_id"`
		Tag    string       `json:"tag"`
	}

	ParamUserTagType struct {
		UserID utils.UInt64 `json:"user_id"`
		Tag    string       `json:"tag"`
		Type   EntityType   `json:"entity_type"`
	}
)

var (
	CKEntity2Tags          = cache.MustNewCacheKey[utils.UInt64]("entity:{entity_id}:tags", CacheExpiration)
	CKUser2Tags            = cache.MustNewCacheKey[utils.UInt64]("entity:{user_id}:tags", CacheExpiration)
	CKUserTag2Entities     = cache.MustNewCacheKey[ParamUserTag]("user:{user_id}:tag:{tag}:entities", CacheExpiration)
	CKUserTagType2Entities = cache.MustNewCacheKey[ParamUserTagType]("user:{user_id}:tag:{tag_id}:type:{entity_type}:entities", CacheExpiration)
	CKTag2Users            = cache.MustNewCacheKey[string]("tag:{tag}:users", CacheExpiration)
)

const (
	EntityTypeItem    EntityType = 1
	EntityTypeBook    EntityType = 2
	EntityTypeDungeon EntityType = 3
)

func (UserTag) TableName() string {
	return "user_tags"
}

// GetTagsByEntity retrieves tags associated with a given entity ID.
func GetTagsByEntity(ctx context.Context, tx *gorm.DB, entityID utils.UInt64) ([]string, error) {
	cacheKey := CKEntity2Tags.MustBuild(entityID)
	cachedTags, err := cache.Set().GetAll(ctx, cacheKey)
	if err == nil && len(cachedTags) > 0 && cachedTags[0] != CacheInvalidationFlag {
		return cachedTags, nil
	}

	var tags []string
	if err = tx.Model(&UserTag{}).Where("entity_id = ? AND deleted_at IS NULL", entityID).Pluck("tag", &tags).Error; err != nil {
		return nil, irr.Wrap(err, "failed to get tags by entity")
	}

	if err = cache.Set().Insert(ctx, cacheKey, CacheExpiration, tags...); err != nil {
		return nil, err
	}

	return tags, nil
}

// AddEntityTags adds tags to an entity for a user.
func AddEntityTags(ctx context.Context, tx *gorm.DB, userID utils.UInt64, entityType EntityType, entityID utils.UInt64, tags ...string) error {
	uts := make([]UserTag, 0, len(tags))
	for _, tag := range tags {
		userTag := UserTag{
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

	if err := invalidateCacheForTag(ctx, tags...); err != nil {
		return err
	}
	return nil
}

// RemoveEntityTags removes tags from an entity for a user.
func RemoveEntityTags(ctx context.Context, tx *gorm.DB, userID utils.UInt64, entityID utils.UInt64, tags ...string) error {
	if err := tx.Model(&UserTag{}).Where("user_id = ? AND tag IN ? AND entity_id = ?", userID, tags, entityID).Update("deleted_at", time.Now()).Error; err != nil {
		return irr.Wrap(err, "failed to remove entity tag")
	}

	// Invalidate relevant caches
	if err := invalidateCacheForEntity(ctx, entityID); err != nil {
		return err
	}
	if err := invalidateCacheForTag(ctx, tags...); err != nil {
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
	cachedEntities, err := cache.Set().GetAllUInt64s(ctx, cacheKey)
	if err == nil && len(cachedEntities) > 0 && cachedEntities[0] != utils.UInt64(0) {
		return cachedEntities, nil
	}

	var entityIDs []utils.UInt64
	if err = tx.Model(&UserTag{}).Where("user_id = ? AND tag = ? AND deleted_at IS NULL", userID, tag).Pluck("entity_id", &entityIDs).Error; err != nil {
		return nil, irr.Wrap(err, "failed to get entities by user and tag")
	}

	if err = cache.Set().InsertUInt64s(ctx, cacheKey, CacheExpiration, entityIDs...); err != nil {
		return nil, err
	}

	return entityIDs, nil
}

// GetEntitiesOfType retrieves entity IDs associated with a given user ID, tag, and entity type.
func GetEntitiesOfType(ctx context.Context, tx *gorm.DB, userID utils.UInt64, tag string, entityType EntityType) ([]utils.UInt64, error) {
	cacheKey := CKUserTagType2Entities.MustBuild(ParamUserTagType{UserID: userID, Tag: tag, Type: entityType})
	cachedEntities, err := cache.Set().GetAllUInt64s(ctx, cacheKey)
	if err == nil && len(cachedEntities) > 0 && cachedEntities[0] != utils.UInt64(0) {
		return cachedEntities, nil
	}

	var entityIDs []utils.UInt64
	if err = tx.Model(&UserTag{}).Where("user_id = ? AND tag = ? AND entity_type = ? AND deleted_at IS NULL", userID, tag, entityType).Pluck("entity_id", &entityIDs).Error; err != nil {
		return nil, irr.Wrap(err, "failed to get entities by user, tag, and entity type")
	}

	if err = cache.Set().InsertUInt64s(ctx, cacheKey, CacheExpiration, entityIDs...); err != nil {
		return nil, err
	}

	return entityIDs, nil
}

// GetTagsByUser retrieves tags associated with a given user ID.
func GetTagsByUser(ctx context.Context, tx *gorm.DB, userID utils.UInt64) ([]string, error) {
	cacheKey := CKUser2Tags.MustBuild(userID)
	cachedTags, err := cache.Set().GetAll(ctx, cacheKey)
	if err == nil && len(cachedTags) > 0 && cachedTags[0] != CacheInvalidationFlag {
		return cachedTags, nil
	}

	var tags []string
	if err := tx.Model(&UserTag{}).Where("user_id = ? AND deleted_at IS NULL").Distinct().Pluck("tag", &tags).Error; err != nil {
		return nil, irr.Wrap(err, "failed to get tags by user")
	}

	if err := cache.Set().Insert(ctx, cacheKey, CacheExpiration, tags...); err != nil {
		return nil, err
	}

	return tags, nil
}

// GetUsersByTag retrieves user IDs associated with a given tag.
func GetUsersByTag(ctx context.Context, tx *gorm.DB, tag string) ([]utils.UInt64, error) {
	cacheKey := CKTag2Users.MustBuild(tag)
	cachedUsers, err := cache.Set().GetAllUInt64s(ctx, cacheKey)
	if err == nil && len(cachedUsers) > 0 && cachedUsers[0] != utils.UInt64(0) {
		return cachedUsers, nil
	}

	var userIDs []utils.UInt64
	if err := tx.Model(&UserTag{}).Where("tag = ? AND deleted_at IS NULL").Distinct().Pluck("user_id", &userIDs).Error; err != nil {
		return nil, irr.Wrap(err, "failed to get users by tag")
	}

	if err := cache.Set().InsertUInt64s(ctx, cacheKey, CacheExpiration, userIDs...); err != nil {
		return nil, err
	}

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
	var existingTags []UserTag
	if err := tx.Where("user_id = ? AND tag = ? AND entity_id IN ?", userID, newTag, entityIDs).Find(&existingTags).Error; err != nil {
		return irr.Wrap(err, "failed to find existing tags")
	}

	existingTagsMap := make(map[utils.UInt64]UserTag)
	for _, tag := range existingTags {
		existingTagsMap[tag.EntityID] = tag
	}

	var tagsToCreate []UserTag
	var tagsToUpdate []UserTag
	for _, entityID := range entityIDs {
		if _, exists := existingTagsMap[entityID]; !exists {
			tagsToCreate = append(tagsToCreate, UserTag{
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
	if err = tx.Model(&UserTag{}).Where(
		"user_id = ? AND tag = ? AND entity_id IN ?",
		userID, newTag, typer.SliceMap(tagsToUpdate, func(from UserTag) utils.UInt64 {
			return from.EntityID
		})).
		Updates(UserTag{DeletedAt: gorm.DeletedAt{}}).Error; err != nil {
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
	if err := tx.Model(&UserTag{}).Where("tag = ?", oldTag).Update("deleted_at", time.Now()).Error; err != nil {
		return irr.Wrap(err, "failed to soft delete old tag")
	}

	// 清除缓存
	if err := invalidateCacheForTag(ctx, oldTag); err != nil {
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
func invalidateCacheForTag(ctx context.Context, tags ...string) error {
	for _, tag := range tags {
		userIDs, err := GetUsersByTag(ctx, nil, tag)
		if err != nil {
			return err
		}

		for _, userID := range userIDs {
			if err = cache.Set().Clear(ctx,
				CKUser2Tags.MustBuild(userID),
				MaxRetryAttempts); err != nil {
				return err
			}
			if err = cache.Set().Clear(ctx,
				CKUserTag2Entities.MustBuild(ParamUserTag{UserID: userID, Tag: tag}),
				MaxRetryAttempts); err != nil {
				return err
			}
			if err = cache.Set().Clear(ctx,
				CKUserTagType2Entities.MustBuild(ParamUserTagType{UserID: userID, Tag: tag, Type: EntityTypeItem}),
				MaxRetryAttempts); err != nil {
				return err
			}
			if err = cache.Set().Clear(ctx,
				CKUserTagType2Entities.MustBuild(ParamUserTagType{UserID: userID, Tag: tag, Type: EntityTypeBook}),
				MaxRetryAttempts); err != nil {
				return err
			}
			if err = cache.Set().Clear(ctx,
				CKUserTagType2Entities.MustBuild(ParamUserTagType{UserID: userID, Tag: tag, Type: EntityTypeDungeon}),
				MaxRetryAttempts); err != nil {
				return err
			}
		}
		if err = cache.Set().Clear(ctx, CKTag2Users.MustBuild(tag), MaxRetryAttempts); err != nil {
			return err
		}
	}
	return nil
}

// invalidateCacheForEntity invalidates caches related to a given entity ID.
func invalidateCacheForEntity(ctx context.Context, entityID utils.UInt64) error {
	return cache.Set().Clear(ctx, CKEntity2Tags.MustBuild(entityID), MaxRetryAttempts)
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
