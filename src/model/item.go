package model

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/khicago/got/util/typer"

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

const (
	TyItemFlashCard      = "flash_card"
	TyItemMultipleChoice = "multiple_choice"
	TyItemCompletion     = "completion"
)

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

// BeforeDelete is a GORM hook that is called before deleting an item.
func (i *Item) BeforeDelete(tx *gorm.DB) (err error) {
	log := wlog.Common("BeforeDeleteItem")
	log.Infof("Deleting associations for item ID %d", i.ID)

	// todo: refine this
	if err = tx.Where("item_id = ?", i.ID).Delete(&BookItem{}).Error; err != nil {
		return irr.Wrap(err, "failed to delete item books")
	}

	return nil
}

func FindItems(ctx context.Context, tx *gorm.DB, itemIDs []utils.UInt64) ([]Item, error) {
	items := make([]Item, 0, len(itemIDs))
	if err := tx.Where("id in ?", itemIDs).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// GetItemIDsOfBooks 获取 book 关联的 items
func GetItemIDsOfBooks(tx *gorm.DB, bookIDs []utils.UInt64) (itemBookMap map[utils.UInt64]utils.UInt64, err error) {
	var bookItems []BookItem
	if err = tx.Where("book_id IN ?", bookIDs).Find(&bookItems).Error; err != nil {
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
	return GetItemsByID(tx, ids)
}

func GetItemsByID(tx *gorm.DB, itemIDs []utils.UInt64) (items []*Item, err error) {
	if err = tx.Where("id in ?", itemIDs).Find(&items).Error; err != nil {
		return nil, irr.Wrap(err, "get items from ids failed")
	}
	return items, nil
}

// GetItemIDsOfTags 获取 tag 关联的 items
func GetItemIDsOfTags(ctx context.Context, userID utils.UInt64, tags []string) (map[utils.UInt64]string, error) {
	itemTagMap := make(map[utils.UInt64]string)
	for _, tag := range tags {

		itemsIDs, err := TagModel().GetEntities(ctx, userID, tag, typer.Ptr(EntityTypeItem))
		if err != nil {
			return nil, err
		}
		for _, itemID := range itemsIDs {
			itemTagMap[itemID] = tag
		}
	}
	return itemTagMap, nil
}

func CreateItems(ctx context.Context, tx *gorm.DB, userID utils.UInt64, items []*Item, itemTagRef map[*Item][]string) error {
	log := wlog.ByCtx(ctx, "model.save_items")

	for _, item := range items {
		if err := tx.Create(item).Error; err != nil {
			return err
		}
		if itemTagRef == nil {
			continue
		}
		if tags, ok := itemTagRef[item]; ok && len(tags) > 0 {
			tags = typer.SliceFilter(tags, func(s string) bool {
				return strings.TrimSpace(s) != ""
			})

			if err := AddEntityTags(ctx, tx, userID, EntityTypeItem, item.ID, tags...); err != nil {
				return irr.Wrap(err, "update item tags failed")
			}
		}
	}
	log.Infof("Successfully saved %d items", len(items))
	return nil
}
