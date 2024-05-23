package model

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/bagaking/goulp/wlog"
	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/khicago/got/util/typer"
	"github.com/khicago/irr"
	"gorm.io/gorm"
)

// Tag represents a tag entity with unique name.
type Tag struct {
	ID   utils.UInt64 `gorm:"primaryKey;autoIncrement:false"`
	Name string       `gorm:"unique;not null"`
}

// TagRefType defines the type of tag reference, such as book or item.
type TagRefType int

const (
	// BookTagRef indicates a reference to a book.
	BookTagRef TagRefType = iota
	// ItemTagRef indicates a reference to an item.
	ItemTagRef
)

var (
	// ErrInvalidTagName indicates that the tag name is invalid.
	ErrInvalidTagName = errors.New("invalid tag name")

	// Default values for batch processing
	defaultBatchSize     = 500
	defaultMaxGoroutines = 10
	defaultTimeout       = 20 * time.Second
)

// FindOrUpdateTagByName finds a tag by name, or creates it if not found.
func FindOrUpdateTagByName(ctx context.Context, tx *gorm.DB, tagName string, id utils.UInt64) (*Tag, error) {
	log := wlog.ByCtx(ctx, "FindOrUpdateTagByName")
	tagName = strings.TrimSpace(tagName)
	if tagName == "" {
		return nil, ErrInvalidTagName
	}

	// 查找或新建 Tag
	tag := &Tag{}
	err := tx.Where(&Tag{Name: tagName}).First(tag).Error
	if err == nil {
		return tag, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, irr.Wrap(err, "find tag failed")
	}

	tag = &Tag{ID: id, Name: tagName}
	if err = tx.Create(tag).Error; err != nil {
		return nil, irr.Wrap(err, "create tag failed")
	}

	log.Infof("new tag %v created, id_set= %d", tag, id)
	return tag, nil
}

// UpdateBookTagsRef updates the tags associated with a book.
func UpdateBookTagsRef(ctx context.Context, tx *gorm.DB, bookID utils.UInt64, tags []string) error {
	return updateTagsRef(ctx, tx, BookTagRef, bookID, tags)
}

// UpdateItemTagsRef updates the tags associated with an item.
func UpdateItemTagsRef(ctx context.Context, tx *gorm.DB, itemID utils.UInt64, tags []string) error {
	return updateTagsRef(ctx, tx, ItemTagRef, itemID, tags)
}

// GetBookTagNames updates the tags associated with a book.
func GetBookTagNames(ctx context.Context, tx *gorm.DB, bookID utils.UInt64) ([]string, error) {
	return getTagsByEntityID(ctx, tx, BookTagRef, bookID)
}

// GetItemTagNames updates the tags associated with an item.
func GetItemTagNames(ctx context.Context, tx *gorm.DB, itemID utils.UInt64) ([]string, error) {
	return getTagsByEntityID(ctx, tx, ItemTagRef, itemID)
}

// getTagsByEntityID 根据实体ID获取关联的标签列表.
func getTagsByEntityID(ctx context.Context, tx *gorm.DB, entityType TagRefType, entityID utils.UInt64) ([]string, error) {
	existingTagIDs, err := fetchTagIDs(tx, entityType, entityID)
	if err != nil {
		return nil, irr.Wrap(err, "failed to fetch existing tag IDs")
	}

	var tags []Tag
	if err = tx.Where("id IN ?", existingTagIDs).Find(&tags).Error; err != nil {
		return nil, irr.Wrap(err, "failed to get tags")
	}

	tagNames := make([]string, len(tags))
	for i, tag := range tags {
		tagNames[i] = tag.Name
	}

	return tagNames, nil
}

// updateTagsRef updates the tags associated with an entity (book or item).
func updateTagsRef(ctx context.Context, tx *gorm.DB, entityType TagRefType, entityID utils.UInt64, tags []string) error {
	log := wlog.ByCtx(ctx, "updateTagsRef")
	tagIDs, err := utils.MGenIDU64(ctx, len(tags))
	if err != nil {
		return irr.Wrap(err, "failed to generate IDs for tags")
	}

	// Fetch and map existing tags
	existingTagMap, err := getExistingTagMap(ctx, tx, entityType, entityID)
	if err != nil {
		return err
	}

	var tagsToAssociate []any
	for i, tagName := range tags {
		tagID, exists := existingTagMap[tagName]
		if !exists {
			tag, err := FindOrUpdateTagByName(ctx, tx, tagName, tagIDs[i])
			if err != nil {
				if errors.Is(err, ErrInvalidTagName) {
					continue
				}
				return irr.Wrap(err, "upsert tag failed")
			}
			tagID = tag.ID
		}

		tagsToAssociate = append(tagsToAssociate, createTagAssoc(entityType, entityID, tagID))
	}

	// Identify and remove obsolete tags
	if err = removeObsoleteTags(ctx, tx, entityType, entityID, existingTagMap, tags); err != nil {
		return err
	}

	// Create new tag associations
	if err = tx.Create(&tagsToAssociate).Error; err != nil {
		return irr.Wrap(err, "failed to associate new tags with entity")
	}

	log.Infof("Tags updated successfully for entity %d of type %d", entityID, entityType)
	return nil
}

// getExistingTagMap retrieves existing tags and constructs a map for quick lookup.
func getExistingTagMap(ctx context.Context, tx *gorm.DB, entityType TagRefType, entityID utils.UInt64) (map[string]utils.UInt64, error) {
	existingTagIDs, err := fetchTagIDs(tx, entityType, entityID)
	if err != nil {
		return nil, irr.Wrap(err, "failed to fetch existing tag IDs")
	}

	existingTags, err := fetchTagsInBatchesParallel(ctx, tx, existingTagIDs)
	if err != nil {
		return nil, irr.Wrap(err, "failed to fetch existing tags")
	}

	existingTagMap := make(map[string]utils.UInt64)
	for _, tag := range existingTags {
		existingTagMap[tag.Name] = tag.ID
	}

	return existingTagMap, nil
}

// createTagAssoc creates a tag association struct based on entity type.
func createTagAssoc(entityType TagRefType, entityID, tagID utils.UInt64) interface{} {
	switch entityType {
	case BookTagRef:
		return BookTag{BookID: entityID, TagID: tagID}
	case ItemTagRef:
		return ItemTag{ItemID: entityID, TagID: tagID}
	default:
		return nil
	}
}

// removeObsoleteTags identifies and removes tags that are no longer associated with the entity.
func removeObsoleteTags(ctx context.Context, tx *gorm.DB, entityType TagRefType, entityID utils.UInt64, existingTagMap map[string]utils.UInt64, tags []string) error {
	var tagsToDelete []utils.UInt64
	for name, id := range existingTagMap {
		if !typer.SliceContains(tags, name) {
			tagsToDelete = append(tagsToDelete, id)
		}
	}

	if len(tagsToDelete) == 0 {
		return nil
	}

	var err error
	switch entityType {
	case BookTagRef:
		err = tx.Where("book_id = ? AND tag_id IN (?)", entityID, tagsToDelete).Delete(&BookTag{}).Error
	case ItemTagRef:
		err = tx.Where("item_id = ? AND tag_id IN (?)", entityID, tagsToDelete).Delete(&ItemTag{}).Error
	}
	return irr.Wrap(err, "failed to delete obsolete tags")
}

// fetchTagIDs fetches tag IDs associated with an entity.
func fetchTagIDs(tx *gorm.DB, entityType TagRefType, entityID utils.UInt64) ([]utils.UInt64, error) {
	var tagIDs []utils.UInt64
	switch entityType {
	case BookTagRef:
		tx = tx.Model(&BookTag{}).Select("tag_id").Where("book_id = ?", entityID)
	case ItemTagRef:
		tx = tx.Model(&ItemTag{}).Select("tag_id").Where("item_id = ?", entityID)
	default:
		return nil, errors.New("unsupported entity type")
	}

	rows, err := tx.Rows()
	if err != nil {
		return nil, irr.Wrap(err, "failed to fetch tag ids")
	}
	defer rows.Close()

	for rows.Next() {
		var id utils.UInt64
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		tagIDs = append(tagIDs, id)
	}

	return tagIDs, nil
}

// fetchTagsInBatchesParallel fetches tags in parallel batches.
func fetchTagsInBatchesParallel(ctx context.Context, tx *gorm.DB, tagIDs []utils.UInt64) ([]Tag, error) {
	var mutex sync.Mutex
	var tags []Tag

	// Processor function to fetch a batch of tags
	processor := func(ctx context.Context, start, end int) error {
		var batch []Tag
		// Fetch a batch of tags within the specified range from the database
		if err := tx.Model(&Tag{}).Where("id IN (?)", tagIDs[start:end]).Find(&batch).Error; err != nil {
			return err
		}
		// Lock the mutex before appending to the shared slice
		mutex.Lock()
		tags = append(tags, batch...)
		mutex.Unlock()
		return nil
	}

	// Execute the batch processing function with the specified options
	if err := utils.ParallelBatchProcess(ctx, len(tagIDs), processor,
		utils.BatchWithSize(defaultBatchSize),
		utils.BatchWithMaxGoroutines(defaultMaxGoroutines),
		utils.BatchWithTimeout(defaultTimeout),
	); err != nil {
		return nil, err
	}

	return tags, nil
}
