package model

import (
	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/def"
	"github.com/khicago/got/util/typer"
	"github.com/khicago/irr"
	"gorm.io/gorm"
)

type (
	// UserMonster - Item 对特定用户的属性
	UserMonster struct {
		UserID utils.UInt64
		ItemID utils.UInt64

		Familiarity utils.Percentage `gorm:"default:0"` // 熟练度，范围为0-100，默认值为0

		monster *Item
	}

	DungeonMonster struct {
		DungeonID utils.UInt64
		ItemID    utils.UInt64

		SourceType MonsterSource
		SourceID   utils.UInt64

		Visibility utils.Percentage `gorm:"default:0"` // 显影程度，根据复习次数变化

		monster *Item
	}

	MonsterSource uint8
)

const (
	MonsterSourceItem MonsterSource = iota + 1
	MonsterSourceBook
	MonsterSourceTag
)

func (dm *DungeonMonster) Monster() *Item {
	return dm.monster
}

func (d *Dungeon) AddMonster(tx *gorm.DB, source MonsterSource, sourceEntityIDs []utils.UInt64) error {
	// Validate the existence of resources
	if err := validateExistence(tx, source, sourceEntityIDs); err != nil {
		return irr.Wrap(err, "add monster to dungeon failed")
	}

	if err := createDungeonRef(tx, source, d.ID, sourceEntityIDs, d.Type == def.DungeonTypeCampaign); err != nil {
		return irr.Wrap(err, "add monster to dungeon failed")
	}

	return nil
}

func validateExistence(tx *gorm.DB, source MonsterSource, resourceIDs []utils.UInt64) error {
	var count int64
	switch source {
	case MonsterSourceBook:
		if err := tx.Model(&Book{}).Where("id IN ?", resourceIDs).Count(&count).Error; err != nil {
			return err
		}
	case MonsterSourceItem:
		if err := tx.Model(&Item{}).Where("id IN ?", resourceIDs).Count(&count).Error; err != nil {
			return err
		}
	case MonsterSourceTag:
		if err := tx.Model(&Tag{}).Where("id IN ?", resourceIDs).Count(&count).Error; err != nil {
			return err
		}
	default:
		return irr.Error("unknown resource type: %v", source)
	}

	if count != int64(len(resourceIDs)) {
		return irr.Error("some resources not found")
	}

	return nil
}

func createDungeonRef(tx *gorm.DB, source MonsterSource, dungeonID utils.UInt64, sourceEntityIDs []utils.UInt64, loadAssociationMonster bool) error {
	for _, id := range sourceEntityIDs {
		switch source {
		case MonsterSourceItem:
			if err := createDungeonMonster(tx, dungeonID, id, source, id); err != nil {
				return err
			}
		case MonsterSourceBook:
			if err := createDungeonBookRecord(tx, dungeonID, id); err != nil {
				return err
			}
			if loadAssociationMonster {
				if err := createMonstersForBook(tx, dungeonID, id); err != nil {
					return err
				}
			}
		case MonsterSourceTag:
			if err := createDungeonTagRecord(tx, dungeonID, id); err != nil {
				return err
			}
			if loadAssociationMonster {
				if err := createMonstersForTag(tx, dungeonID, id); err != nil {
					return err
				}
			}
		default:
			return irr.Error("unknown resource type: %v", source)
		}
	}
	return nil
}

func createDungeonBookRecord(tx *gorm.DB, dungeonID, bookID utils.UInt64) error {
	dungeonBook := DungeonBook{
		DungeonID: dungeonID,
		BookID:    bookID,
	}
	if err := tx.Where(&dungeonBook).FirstOrCreate(&dungeonBook).Error; err != nil {
		return err
	}
	return nil
}

func createDungeonTagRecord(tx *gorm.DB, dungeonID, tagID utils.UInt64) error {
	dungeonTag := DungeonTag{
		DungeonID: dungeonID,
		TagID:     tagID,
	}
	if err := tx.Where(&dungeonTag).FirstOrCreate(&dungeonTag).Error; err != nil {
		return err
	}
	return nil
}

func createDungeonMonster(tx *gorm.DB, dungeonID, itemID utils.UInt64, source MonsterSource, sourceEntityID utils.UInt64) error {
	dungeonMonster := DungeonMonster{
		DungeonID:  dungeonID,
		ItemID:     itemID,
		SourceType: source,
		SourceID:   sourceEntityID,
		Visibility: 0,
	}
	if err := tx.Where("dungeon_id = ? AND item_id = ?", dungeonID, itemID).FirstOrCreate(&dungeonMonster).Error; err != nil {
		return err
	}
	return nil
}

func createMonstersForBook(tx *gorm.DB, dungeonID, bookID utils.UInt64) error {
	var bookItems []BookItem
	if err := tx.Where("book_id = ?", bookID).Find(&bookItems).Error; err != nil {
		return err
	}
	for _, bookItem := range bookItems {
		if err := createDungeonMonster(tx, dungeonID, bookItem.ItemID, MonsterSourceItem, bookID); err != nil {
			return err
		}
	}
	return nil
}

func createMonstersForTag(tx *gorm.DB, dungeonID, tagID utils.UInt64) error {
	var tagItems []ItemTag
	if err := tx.Where("tag_id = ?", tagID).Find(&tagItems).Error; err != nil {
		return err
	}
	for _, tagItem := range tagItems {
		if err := createDungeonMonster(tx, dungeonID, tagItem.ItemID, MonsterSourceItem, tagID); err != nil {
			return err
		}
	}
	return nil
}

// GetMonsters - 获取当前 Dungeon 的已创建的 DungeonMonster
func (d *Dungeon) GetMonsters(tx *gorm.DB, sortBy string, offset, limit int) ([]DungeonMonster, error) {
	var monsters []DungeonMonster
	query := tx.Where("dungeon_id = ?", d.ID).Offset(offset).Limit(limit)

	// todo: 要去查 item
	//switch sortBy {
	//case "familiarity":
	//	query = query.Order("familiarity ASC, item_id ASC")
	//case "difficulty":
	//	query = query.Order("difficulty ASC, item_id ASC")
	//case "importance":
	//	query = query.Order("importance ASC, item_id ASC")
	//default:
	//	query = query.Order("item_id ASC")
	//}

	if err := query.Find(&monsters).Error; err != nil {
		return nil, err
	}

	return monsters, nil
}

// GetMonstersWithAssociations - 获取当前 Dungeon 的 DungeonMonster 及其关联的 Items, Books, TagNames
func (d *Dungeon) GetMonstersWithAssociations(tx *gorm.DB, sortBy string, offset, limit int) ([]DungeonMonster, error) {
	// 获取关联的 Items, Books, TagNames
	books, items, tags, err := GetDungeonAssociations(tx, d.ID)
	if err != nil {
		return nil, err
	}

	// 直接获取 items
	itemSourceMap := make(map[utils.UInt64]utils.UInt64)
	for _, itemID := range items {
		itemSourceMap[itemID] = itemID
	}

	// 获取 book 关联的 items
	bookItemMap, err := GetItemsOfBooks(tx, books)
	if err != nil {
		return nil, err
	}
	for itemID := range bookItemMap {
		if bookID, exists := itemSourceMap[itemID]; exists {
			itemSourceMap[itemID] = bookID
		}
	}

	// 获取 tag 关联的 items
	tagItemMap, err := GetItemsOfTags(tx, tags)
	if err != nil {
		return nil, err
	}
	for itemID := range tagItemMap {
		if tagID, exists := itemSourceMap[itemID]; exists {
			itemSourceMap[itemID] = tagID
		}
	}

	// 批量获取所有 item 的详细信息
	itemIDs := typer.Keys(itemSourceMap)

	// 获取所有 item 的详细信息并排序分页
	txItems := tx.Table("items").Where("id IN (?)", itemIDs)
	switch sortBy {
	case "familiarity":
		txItems = txItems.Order("familiarity ASC, id ASC")
	case "difficulty":
		txItems = txItems.Order("difficulty ASC, id ASC")
	case "importance":
		txItems = txItems.Order("importance ASC, id ASC")
	default:
		txItems = txItems.Order("id ASC")
	}

	var itemsList []*Item
	if err = txItems.Offset(offset).Limit(limit).Find(&itemsList).Error; err != nil {
		return nil, err
	}

	// 转换 itemsList 为 monsters slice，注意这个方法并没有查询 dungeon_monster 表
	monsters := make([]DungeonMonster, 0, len(itemsList))
	for _, item := range itemsList {
		monster := DungeonMonster{
			ItemID:     item.ID,
			DungeonID:  d.ID,
			SourceType: MonsterSourceItem,
			SourceID:   itemSourceMap[item.ID],
			monster:    item,
		}
		if _, ok := tagItemMap[item.ID]; ok {
			monster.SourceType = MonsterSourceTag
		}
		if _, ok := bookItemMap[item.ID]; ok {
			monster.SourceType = MonsterSourceBook
		}
		monsters = append(monsters, monster)
	}

	return monsters, nil
}
