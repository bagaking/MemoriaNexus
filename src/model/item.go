package model

import (
	"errors"
	"time"

	"github.com/bagaking/memorianexus/internal/utils"

	"github.com/bagaking/memorianexus/src/def"

	"github.com/bagaking/goulp/wlog"
	"github.com/khicago/irr"
	"gorm.io/gorm"
)

type Item struct {
	ID utils.UInt64 `gorm:"primaryKey;autoIncrement:true" json:"id"`

	CreatorID utils.UInt64 `gorm:"not null" json:"creator_id,string"`

	Type    string
	Content string

	// todo: 示意，后续应该放到 user、item 关联的表中去
	Difficulty def.DifficultyLevel `gorm:"default:0x01"` // 难度，默认值为 NoviceNormal (0x01)
	Importance def.ImportanceLevel `gorm:"default:0x01"` // 重要程度，默认值为 DomainGeneral (0x01)

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt
}

func (i *Item) TableName() string {
	return "items"
}

// BeforeCreate 钩子
func (i *Item) BeforeCreate(tx *gorm.DB) (err error) {
	// 确保UserID不为0
	if i.ID <= 0 {
		return errors.New("user UInt64 must be larger than zero")
	}
	return
}

type ItemTag struct {
	ItemID utils.UInt64 `gorm:"primaryKey"`
	TagID  utils.UInt64 `gorm:"primaryKey"`
}

func (ItemTag) Associate(itemID, tagID utils.UInt64) ITagAssociate {
	return ItemTag{ItemID: itemID, TagID: tagID}
}

func (ItemTag) Type() TagRefType {
	return ItemTagRef
}

var _ ITagAssociate = &ItemTag{}

func (b ItemTag) TableName() string {
	return "item_tags"
}

// BeforeDelete is a GORM hook that is called before deleting an item.
func (i *Item) BeforeDelete(tx *gorm.DB) (err error) {
	log := wlog.Common("BeforeDeleteItem")
	log.Infof("Deleting associations for item ID %d", i.ID)

	// todo: refine this

	if err = tx.Where("item_id = ?", i.ID).Delete(&ItemTag{}).Error; err != nil {
		return irr.Wrap(err, "failed to delete item tags")
	}

	if err = tx.Where("item_id = ?", i.ID).Delete(&BookItem{}).Error; err != nil {
		return irr.Wrap(err, "failed to delete item books")
	}

	return nil
}

// GetItemsOfBooks 获取 book 关联的 items
func GetItemsOfBooks(tx *gorm.DB, bookIDs []utils.UInt64) (itemBookMap map[utils.UInt64]utils.UInt64, err error) {
	var bookItems []BookItem
	if err = tx.Where("book_id IN (?)", bookIDs).Find(&bookItems).Error; err != nil {
		return nil, err
	}
	itemBookMap = make(map[utils.UInt64]utils.UInt64)
	for _, bookItem := range bookItems {
		itemBookMap[bookItem.ItemID] = bookItem.BookID
	}
	return itemBookMap, nil
}

// GetItemIDsOfBook 获取某个 book 的 items
func GetItemIDsOfBook(tx *gorm.DB, bookID utils.UInt64, offset, limit int) (itemIDs []utils.UInt64, err error) {
	if err = tx.Model(&BookItem{}).Where(
		"book_id = ?", bookID,
	).Offset(offset).Limit(limit).Pluck("item_id", &itemIDs).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	return itemIDs, nil
}

// GetItemsOfBook 获取某个 book 的 items
func GetItemsOfBook(tx *gorm.DB, bookID utils.UInt64, offset, limit int) (items []*Item, err error) {
	ids, err := GetItemIDsOfBook(tx, bookID, offset, limit)
	if err != nil {
		return nil, irr.Wrap(err, "get item ids for book failed")
	}
	if err = tx.Where("id in (?)", ids).Find(&items).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, irr.Wrap(err, "get items from ids failed")
		}
	}
	return items, nil
}

// GetItemsOfTags 获取 tag 关联的 items
func GetItemsOfTags(tx *gorm.DB, tagIDs []utils.UInt64) (map[utils.UInt64]utils.UInt64, error) {
	var tagItems []ItemTag
	if err := tx.Where("tag_id IN (?)", tagIDs).Find(&tagItems).Error; err != nil {
		return nil, err
	}
	itemTagMap := make(map[utils.UInt64]utils.UInt64)
	for _, tagItem := range tagItems {
		itemTagMap[tagItem.ItemID] = tagItem.TagID
	}
	return itemTagMap, nil
}
