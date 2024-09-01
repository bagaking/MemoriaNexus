package model

import (
	"context"
	"time"

	"github.com/adjust/redismq"
	"gorm.io/gorm"

	"github.com/bagaking/goulp/wlog"
	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/pkg/tags"
	"github.com/khicago/got/util/typer"
	"github.com/khicago/irr"
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

	TagRepo struct {
		db *gorm.DB
	}

	TModel struct {
		*tags.TagService[EntityType]
	}
)

func (t TagRepo) GetTag(ctx context.Context, tag string) (*tags.Tag[EntityType], error) {
	var tagModel Tag
	if err := t.db.WithContext(ctx).Where("tag = ?", tag).First(&tagModel).Error; err != nil {
		return nil, irr.Wrap(err, "failed to get tag")
	}

	return &tags.Tag[EntityType]{
		UserID:     tagModel.UserID,
		Tag:        tagModel.Tag,
		EntityID:   tagModel.EntityID,
		EntityType: tagModel.EntityType,
		CreatedAt:  tagModel.CreatedAt,
		UpdatedAt:  tagModel.UpdatedAt,
	}, nil
}

func (t TagRepo) GetTagsByUser(ctx context.Context, userID utils.UInt64) ([]string, error) {
	var tags []string
	if err := t.db.WithContext(ctx).Model(&Tag{}).Where("user_id = ?", userID).Pluck("tag", &tags).Error; err != nil {
		return nil, irr.Wrap(err, "failed to get tags by user")
	}

	return tags, nil
}

func (t TagRepo) GetEntitiesByTag(ctx context.Context, userID utils.UInt64, tag string, entityType EntityType) ([]utils.UInt64, error) {
	var entityIDs []utils.UInt64
	if err := t.db.WithContext(ctx).Model(&Tag{}).Where("user_id = ? AND tag = ? AND entity_type = ?", userID, tag, entityType).Pluck("entity_id", &entityIDs).Error; err != nil {
		return nil, irr.Wrap(err, "failed to get entities by tag")
	}

	return entityIDs, nil
}

func (t TagRepo) AddTags(ctx context.Context, userID utils.UInt64, entityID utils.UInt64, entityType EntityType, tags ...string) error {
	tagModels := make([]Tag, len(tags))
	for i, tag := range tags {
		tagModels[i] = Tag{
			UserID:     userID,
			Tag:        tag,
			EntityID:   entityID,
			EntityType: entityType,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
	}

	if err := t.db.WithContext(ctx).Create(&tagModels).Error; err != nil {
		return irr.Wrap(err, "failed to add tags")
	}

	return nil
}

func (t TagRepo) RemoveTags(ctx context.Context, userID utils.UInt64, entityID utils.UInt64, tags ...string) error {
	if err := t.db.WithContext(ctx).Where("user_id = ? AND entity_id = ? AND tag IN ?", userID, entityID, tags).Delete(&Tag{}).Error; err != nil {
		return irr.Wrap(err, "failed to remove tags")
	}

	return nil
}

func (t TagRepo) RenameTag(ctx context.Context, userID utils.UInt64, oldTag, newTag string) error {
	if err := t.db.WithContext(ctx).Model(&Tag{}).Where("user_id = ? AND tag = ?", userID, oldTag).Update("tag", newTag).Error; err != nil {
		return irr.Wrap(err, "failed to rename tag")
	}

	return nil
}

func (t TagRepo) GetTagsByEntity(ctx context.Context, entityID utils.UInt64) ([]string, error) {
	var tags []string
	if err := t.db.WithContext(ctx).Model(&Tag{}).Where("entity_id = ?", entityID).Pluck("tag", &tags).Error; err != nil {
		return nil, irr.Wrap(err, "failed to get tags by entity")
	}

	return tags, nil
}

func (t TagRepo) GetUsersByTag(ctx context.Context, tag string) ([]utils.UInt64, error) {
	var userIDs []utils.UInt64
	if err := t.db.WithContext(ctx).Model(&Tag{}).Where("tag = ?", tag).Pluck("user_id", &userIDs).Error; err != nil {
		return nil, irr.Wrap(err, "failed to get users by tag")
	}

	return userIDs, nil
}

const (
	EntityTypeItem    EntityType = 1
	EntityTypeBook    EntityType = 2
	EntityTypeDungeon EntityType = 3
)

var tagModel *TModel

func MustInit(ctx context.Context, db *gorm.DB, queue *redismq.Queue) {
	producer := tags.NewRedisMQProducer(queue)
	consumer, err := tags.NewRedisMQConsumer(queue, "tag_consumer")
	if err != nil {
		wlog.ByCtx(ctx, "model.MustInit").WithError(err).Fatalf("failed to create redis mq consumer")
	}

	tagService := tags.NewTagService[EntityType](
		ctx,
		&TagRepo{
			db: db,
		},
		[]EntityType{
			EntityTypeItem,
			EntityTypeBook,
			EntityTypeDungeon,
		},
		producer, consumer,
	)

	tagModel = &TModel{
		TagService: tagService,
	}
}

func TagModel() *TModel {
	return tagModel
}

func (Tag) TableName() string {
	return "tags"
}

// AddEntityTags adds tags to an entity for a user.
func AddEntityTags(ctx context.Context, tx *gorm.DB, userID utils.UInt64, entityType EntityType, entityID utils.UInt64, tags ...string) error {
	if len(tags) == 0 {
		wlog.ByCtx(ctx, "AddEntityTags").WithField("user_id", userID).WithField("entity_id", entityID).
			Warnf("cannot add tags with empty list")
		return nil
	}
	tagModels := make([]Tag, len(tags))
	for i, tag := range tags {
		tagModels[i] = Tag{
			UserID:     userID,
			Tag:        tag,
			EntityID:   entityID,
			EntityType: entityType,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
	}

	if err := tx.Create(&tagModels).Error; err != nil {
		return irr.Wrap(err, "failed to add entity tag")
	}

	// Invalidate relevant caches
	if err := TagModel().InvalidateEntityCache(ctx, entityID, true); err != nil {
		return err
	}
	// all user's tag-related cache should be cleared
	// todo: consider to clear only the cache with the entity type
	if err := invalidateTagsCache(ctx, tags...); err != nil {
		return err
	}
	return nil
}

// RemoveEntityTags removes tags from an entity for a user.
func RemoveEntityTags(ctx context.Context, tx *gorm.DB, userID utils.UInt64, entityID utils.UInt64, tagsToRemove ...string) error {
	if len(tagsToRemove) == 0 {
		wlog.ByCtx(ctx, "RemoveEntityTags").WithField("user_id", userID).WithField("entity_id", entityID).
			Warnf("cannot remove tags with empty list")
		return nil
	}

	if err := tx.Model(&Tag{}).Where("user_id = ? AND tag IN ? AND entity_id = ?", userID, tagsToRemove, entityID).Update("deleted_at", time.Now()).Error; err != nil {
		return irr.Wrap(err, "failed to remove entity tag")
	}

	// Invalidate relevant caches
	if err := TagModel().InvalidateEntityCache(ctx, entityID, true); err != nil {
		return err
	}
	if err := invalidateTagsCache(ctx, tagsToRemove...); err != nil {
		return err
	}

	return nil
}

// FindItemsOfTag returns the items
func FindItemsOfTag(ctx context.Context, tx *gorm.DB, userID utils.UInt64, tag string, pager *utils.Pager) ([]Item, error) {
	entityType := EntityTypeItem
	itemIDs, err := TagModel().GetEntities(ctx, userID, tag, &entityType)
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
	entityType := EntityTypeBook
	bookIDs, err := TagModel().GetEntities(ctx, userID, tag, &entityType)
	if err != nil {
		return nil, irr.Wrap(err, "failed to find items of tag= %v", tag)
	}
	var books []Book
	if err = tx.Where("id IN ?", bookIDs).Offset(pager.Offset).Limit(pager.Limit).Find(&books).Error; err != nil {
		return nil, irr.Wrap(err, "failed to fetch books by id in %v", bookIDs)
	}
	return books, nil
}

// UpdateEntityTagsDiff handles the difference between the tags of an entity and the new tags.
func UpdateEntityTagsDiff(ctx context.Context, tx *gorm.DB, userID utils.UInt64, entityID utils.UInt64, tags []string) error {
	log := wlog.ByCtx(ctx, "UpdateEntityTagsDiff")
	tagsExist, err := TagModel().GetTagsOfEntity(ctx, entityID)
	if err != nil {
		return irr.Wrap(err, "failed to get exist tags")
	}
	if len(tagsExist) == 0 {
		if len(tags) > 0 {
			if err = AddEntityTags(ctx, tx, userID, EntityTypeItem, entityID, tags...); err != nil {
				return irr.Wrap(err, "failed to add tags")
			}
		}
		return nil
	}

	toAdd, toRemove := typer.SliceDiff(tagsExist, tags)
	log.Debugf("diffLists, add= %v, rem= %v, from= %v, to= %v", toAdd, toRemove, tagsExist, tags)
	if len(toAdd) > 0 {
		if err = AddEntityTags(ctx, tx, userID, EntityTypeItem, entityID, toAdd...); err != nil {
			return irr.Wrap(err, "failed to add tags")
		}
	}
	if len(toRemove) > 0 {
		if err = RemoveEntityTags(ctx, tx, userID, entityID, toRemove...); err != nil {
			return irr.Wrap(err, "failed to remove tags")
		}
	}
	return nil
}

// invalidateTagsCache invalidates caches related to a given tag.
func invalidateTagsCache(ctx context.Context, tags ...string) error {
	for _, tag := range tags {
		if err := TagModel().InvalidateTagCache(ctx, tag, true); err != nil {
			return irr.Wrap(err, "failed to invalidate cache for tag")
		}
	}
	return nil
}
