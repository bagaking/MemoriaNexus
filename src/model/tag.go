package model

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm/clause"

	"github.com/bagaking/memorianexus/internal/utils"

	"gorm.io/gorm"

	"github.com/bagaking/goulp/wlog"
	"github.com/khicago/got/util/typer"
	"github.com/khicago/irr"
)

// Tag represents a tag entity with unique name.
type (
	Tag struct {
		ID   utils.UInt64 `gorm:"primaryKey;autoIncrement:false"`
		Name string       `gorm:"unique;not null"`
	}

	ITagAssociate interface {
		Associate(entityID, tagID utils.UInt64) ITagAssociate
		Type() TagRefType
	}
)

// BeforeCreate 钩子
func (t *Tag) BeforeCreate(tx *gorm.DB) (err error) {
	// 确保UserID不为0
	if t.ID <= 0 {
		return errors.New("user UInt64 must be larger than zero")
	}
	return
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

func (ty TagRefType) CreateAssociate(entityID, tagID utils.UInt64) ITagAssociate {
	switch ty {
	case BookTagRef:
		return BookTag{}.Associate(entityID, tagID)
	case ItemTagRef:
		return ItemTag{}.Associate(entityID, tagID)
	}
	return nil
}

// FindTagsByName fetches tag IDs associated with an entity.
func FindTagsByName(ctx context.Context, tx *gorm.DB, names []string) ([]Tag, error) {
	var tags []Tag
	if err := tx.Where("name in ?", names).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// FindTagByName fetches tag IDs associated with an entity.
func FindTagByName(ctx context.Context, tx *gorm.DB, name string) (*Tag, error) {
	var tag Tag
	if err := tx.Where("name = ?", name).First(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

// FindItemsOfTag returns the items
func FindItemsOfTag(ctx context.Context, tx *gorm.DB, tagID utils.UInt64, pager *utils.Pager) ([]Item, error) {
	var itemTags []ItemTag
	if err := tx.Where("tag_id = ?", tagID).Offset(pager.Offset).Limit(pager.Limit).Find(&itemTags).Error; err != nil {
		return nil, irr.Wrap(err, "failed to find item_ids of tag_id= %v", tagID)
	}

	var itemIDs []utils.UInt64
	for _, itemTag := range itemTags {
		itemIDs = append(itemIDs, itemTag.ItemID)
	}

	var items []Item
	if err := tx.Where("id IN ?", itemIDs).Offset(pager.Offset).Limit(pager.Limit).Find(&items).Error; err != nil {
		return nil, irr.Wrap(err, "failed to fetch items by id in %v", itemIDs)
	}
	return items, nil
}

// FindBooksOfTag returns the items
func FindBooksOfTag(ctx context.Context, tx *gorm.DB, tagID utils.UInt64, pager *utils.Pager) ([]Book, error) {
	var bookTags []BookTag
	if err := tx.Where("tag_id = ?", tagID).Offset(pager.Offset).Limit(pager.Limit).Find(&bookTags).Error; err != nil {
		return nil, irr.Wrap(err, "failed to find book_ids of tag_id= %v", tagID)
	}

	var bookIDs []utils.UInt64
	for _, book := range bookTags {
		bookIDs = append(bookIDs, book.BookID)
	}

	var books []Book
	if err := tx.Where("id IN ?", bookIDs).Offset(pager.Offset).Limit(pager.Limit).Find(&books).Error; err != nil {
		return nil, irr.Wrap(err, "failed to fetch books by id in %v", bookIDs)
	}
	return books, nil
}

// FindTagsIDByName fetches tag IDs associated with an entity.
func FindTagsIDByName(tx *gorm.DB, names []string) ([]utils.UInt64, error) {
	var tagIDs []utils.UInt64
	if err := tx.Model(&Tag{}).Where("name in ?", names).Pluck("id", &tagIDs).Error; err != nil {
		return nil, err
	}
	return tagIDs, nil
}

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

	log.Infof("new tag %+v created, id_set= %d", tag, id)
	return tag, nil
}

// UpdateBookTagsRef updates the tags associated with a book.
func UpdateBookTagsRef(ctx context.Context, tx *gorm.DB, bookID utils.UInt64, tags []string) error {
	return updateTagsRef[BookTag](ctx, tx, bookID, tags)
}

// UpdateItemTagsRef updates the tags associated with an item.
func UpdateItemTagsRef(ctx context.Context, tx *gorm.DB, itemID utils.UInt64, tags []string) error {
	return updateTagsRef[ItemTag](ctx, tx, itemID, tags)
}

// GetBookTagNames updates the tags associated with a book.
func GetBookTagNames(ctx context.Context, tx *gorm.DB, bookID utils.UInt64) ([]string, error) {
	return getTagNamesByEntityID(ctx, tx, BookTagRef, bookID)
}

// GetItemTagNames updates the tags associated with an item.
func GetItemTagNames(ctx context.Context, tx *gorm.DB, itemID utils.UInt64) ([]string, error) {
	return getTagNamesByEntityID(ctx, tx, ItemTagRef, itemID)
}

// getTagNamesByEntityID 根据实体ID获取关联的标签列表.
func getTagNamesByEntityID(ctx context.Context, tx *gorm.DB, entityType TagRefType, entityID utils.UInt64) ([]string, error) {
	existingTagIDs, err := fetchTagIDsByEntity(tx, entityType, entityID)
	if err != nil {
		return nil, irr.Wrap(err, "failed to fetch existing tag IDs")
	}

	var tags []Tag
	if err = tx.Where("id IN ?", existingTagIDs).Find(&tags).Error; err != nil {
		return nil, irr.Wrap(err, "failed to get tags")
	}

	tagNames := make([]string, 0, len(tags))
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}

	return tagNames, nil
}

func updateTagsRef[T ITagAssociate](ctx context.Context, tx *gorm.DB, entityID utils.UInt64, tags []string) error {
	log := wlog.ByCtx(ctx, "updateTagsRef")
	tagIDs, err := utils.MGenIDU64(ctx, len(tags))
	if err != nil {
		return irr.Wrap(err, "failed to generate IDs for tags")
	}

	mm := typer.ZeroVal[T]()
	existingTagMap, err := getExistingTagMap(ctx, tx, mm.Type(), entityID)
	if err != nil {
		return err
	}

	var tagsToAssociate []T
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
		tagsToAssociate = append(tagsToAssociate, mm.Associate(entityID, tagID).(T))
	}
	// 删除旧的关联
	if err = removeObsoleteTags(ctx, tx, mm.Type(), entityID, existingTagMap, tags); err != nil {
		return err
	}
	// 插入新的关联
	if len(tagsToAssociate) > 0 {
		if err = tx.Clauses(clause.OnConflict{
			DoNothing: true, // 如果冲突则不做任何操作
		}).Create(&tagsToAssociate).Error; err != nil {
			log.Warnf("failed to associate book tags %#v to tags %#v", tagsToAssociate, tags)
			return irr.Wrap(err, "failed to associate new tags with book")
		}
	}

	log.Infof("TagNames updated successfully for entity %d of type %d", entityID, mm.Type())
	return nil
}

// getExistingTagMap retrieves existing tags and constructs a map for quick lookup.
func getExistingTagMap(ctx context.Context, tx *gorm.DB, entityType TagRefType, entityID utils.UInt64) (map[string]utils.UInt64, error) {
	existingTagIDs, err := fetchTagIDsByEntity(tx, entityType, entityID)
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
		err = tx.Where("book_id = ? AND tag_id IN ?", entityID, tagsToDelete).Delete(&BookTag{}).Error
	case ItemTagRef:
		err = tx.Where("item_id = ? AND tag_id IN ?", entityID, tagsToDelete).Delete(&ItemTag{}).Error
	}
	if err != nil { // todo: irr wrap 考虑保护 err = nil 的情况?
		return irr.Wrap(err, "failed to delete obsolete tags, entity_type= %v, tags_to_delete= %v", entityType, tagsToDelete)
	}
	return nil
}

// fetchTagIDsByEntity fetches tag IDs associated with an entity.
func fetchTagIDsByEntity(tx *gorm.DB, entityType TagRefType, entityID utils.UInt64) ([]utils.UInt64, error) {
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
		if err := tx.Model(&Tag{}).Where("id IN ?", tagIDs[start:end]).Find(&batch).Error; err != nil {
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
