package model

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/bagaking/goulp/wlog"
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

		// 用于 runtime
		PracticeAt     time.Time // 上次复习时间的记录
		NextPracticeAt time.Time // 下次复习时间
		PracticeCount  uint32    // 复习次数 (考虑到可能会有 merge 次数等逻辑，这里先用一个相对大的空间）

		// Gaming
		Visibility utils.Percentage `gorm:"default:0"` // Visibility 显影程度，根据复习次数变化
		Avatar     string           // 头像

		// 以下为宽表内容，为了加速查询
		Familiarity utils.Percentage `gorm:"default:0"` // UserMonster 向 DungeonMonster 单项同步

		Difficulty def.DifficultyLevel `gorm:"default:0x01"` // Item 向 DungeonMonster 单项同步
		Importance def.ImportanceLevel `gorm:"default:0x01"` // Item 向 DungeonMonster 单项同步

		CreatedAt time.Time
	}

	MonsterSource uint8
)

const (
	MonsterSourceItem MonsterSource = 1
	MonsterSourceBook MonsterSource = 2
	MonsterSourceTag  MonsterSource = 3
)

func (ms MonsterSource) String() string {
	switch ms {
	case MonsterSourceItem:
		return "item"
	case MonsterSourceBook:
		return "book"
	case MonsterSourceTag:
		return "tag"
	default:
		return fmt.Sprintf("unsupported_monster_source(%d)", ms)
	}
}

func (d *Dungeon) AddMonster(tx *gorm.DB, source MonsterSource, sourceEntityIDs []utils.UInt64) error {

	// Validate the existence of resources
	if err := validateExistence(tx, source, sourceEntityIDs); err != nil {
		return irr.Track(err, "add monster to dungeon failed, source= %v, ids= %v", source, sourceEntityIDs)
	}

	if source == MonsterSourceItem {
		if err := createMonstersByItemID(tx, d.ID, sourceEntityIDs); err != nil {
			return irr.Track(err, "add monster (from item list) to dungeon failed")
		}
	} else {
		for _, id := range sourceEntityIDs {
			switch source {
			case MonsterSourceBook:
				if err := createDungeonBookRecord(tx, d.ID, id); err != nil {
					return err
				}
				if d.Type == def.DungeonTypeCampaign {
					if err := createMonstersForBook(tx, d.ID, id); err != nil {
						return irr.Track(err, "add monster (from book's ref) to dungeon failed")
					}
				}
			case MonsterSourceTag:
				if err := createDungeonTagRecord(tx, d.ID, id); err != nil {
					return err
				}
				if d.Type == def.DungeonTypeCampaign {
					if err := createMonstersForTag(tx, d.ID, id); err != nil {
						return irr.Track(err, "add monster (from tag's ref) to dungeon failed")
					}
				}
			default:
				return irr.Trace("add monster failed, unknown resource type: %v", source)
			}
		}
	}
	return nil
}

func validateExistence(tx *gorm.DB, source MonsterSource, resourceIDs []utils.UInt64) error {
	var count int64
	switch source {
	case MonsterSourceItem:
		if err := tx.Model(&Item{}).Where("id IN ( ? )", resourceIDs).Count(&count).Error; err != nil {
			return irr.Track(err, "find items in ids failed, ids=%v", resourceIDs)
		}
	case MonsterSourceBook:
		if err := tx.Model(&Book{}).Where("id IN ( ? )", resourceIDs).Count(&count).Error; err != nil {
			return irr.Track(err, "find books in ids failed, ids=%v", resourceIDs)
		}
	case MonsterSourceTag:
		if err := tx.Model(&Tag{}).Where("id IN ( ? )", resourceIDs).Count(&count).Error; err != nil {
			return irr.Track(err, "find tags in ids failed, ids=%v", resourceIDs)
		}
	default:
		return irr.Trace("validate failed, unknown resource type: %v", source)
	}

	if count != int64(len(resourceIDs)) {
		return irr.Error("some resources not found")
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

func createDungeonMonster(tx *gorm.DB, dungeonID utils.UInt64, item Item, source MonsterSource, sourceEntityID utils.UInt64) error {
	dungeonMonster := DungeonMonster{
		DungeonID: dungeonID,
		ItemID:    item.ID,

		SourceType: source,
		SourceID:   sourceEntityID,

		// 用于 runtime
		Visibility:     0,
		PracticeCount:  0,
		PracticeAt:     time.Now(),
		NextPracticeAt: time.Now(),

		// 以下为宽表内容，为了加速查询 todo: 设计更新机制
		Familiarity: utils.Percentage(0),
		Difficulty:  item.Difficulty,
		Importance:  item.Importance,

		CreatedAt: time.Now(),
	}
	if err := tx.Where("dungeon_id = ? AND item_id = ?", dungeonID, item.ID).FirstOrCreate(&dungeonMonster).Error; err != nil {
		return err
	}
	return nil
}

func createMonstersForBook(tx *gorm.DB, dungeonID, bookID utils.UInt64) error {
	var bookItems []BookItem
	if err := tx.Where("book_id = ?", bookID).Find(&bookItems).Error; err != nil {
		return err
	}

	itemIDs := typer.SliceMap(bookItems, func(from BookItem) utils.UInt64 {
		return from.ItemID
	})

	var items []Item
	if err := tx.Where("id in (?)", itemIDs).Find(&items).Error; err != nil {
		return err
	}

	for _, item := range items {
		if err := createDungeonMonster(tx, dungeonID, item, MonsterSourceItem, bookID); err != nil {
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

	itemIDs := typer.SliceMap(tagItems, func(from ItemTag) utils.UInt64 {
		return from.ItemID
	})

	var items []Item
	if err := tx.Where("id in (?)", itemIDs).Find(&items).Error; err != nil {
		return err
	}

	for _, item := range items {
		if err := createDungeonMonster(tx, dungeonID, item, MonsterSourceItem, tagID); err != nil {
			return err
		}
	}
	return nil
}

func createMonstersByItemID(tx *gorm.DB, dungeonID utils.UInt64, itemIDs []utils.UInt64) error {
	var items []Item
	if err := tx.Where("id in (?)", itemIDs).Find(&items).Error; err != nil {
		return err
	}

	for _, item := range items {
		if err := createDungeonMonster(tx, dungeonID, item, MonsterSourceItem, item.ID); err != nil {
			return err
		}
	}
	return nil
}

// GetMonsters retrieves the monsters for the dungeon with sorting and pagination
func (d *Dungeon) GetMonsters(tx *gorm.DB, offset, limit int) ([]DungeonMonster, error) {
	var dungeonMonsters []DungeonMonster

	if err := tx.Where("dungeon_id = ?", d.ID).
		Order("item_id ASC").Offset(offset).Limit(limit).
		Find(&dungeonMonsters).Error; err != nil {
		return nil, err
	}

	return dungeonMonsters, nil
}

// GetMonstersForPractice retrieves the monsters for practice based on the memorization strategy
func (d *Dungeon) GetMonstersForPractice(ctx context.Context, tx *gorm.DB, strategy string, count int) ([]DungeonMonster, error) {
	now := time.Now()
	log := wlog.ByCtx(ctx, "GetMonstersForPractice").
		WithField("time", now).
		WithField("strategy", strategy).
		WithField("dungeon_id", d.ID).
		WithField("count", count)

	log.Infof("start get monsters")
	var query *gorm.DB
	var dungeonMonsters []DungeonMonster
	// 根据复习策略进行查询
	switch strategy {
	case "classic":
		// 经典策略：下次复习时间早于当前时间，熟练度最低，重要程度最高，难度最低
		query = tx.Where("dungeon_id = ? AND next_practice_at < ?", d.ID, now).
			Order("familiarity ASC, importance DESC, difficulty ASC").
			Limit(count).Find(&dungeonMonsters)
	default:
		// 默认策略：下次复习时间早于当前时间，按熟练度排序
		query = tx.Where("dungeon_id = ? AND next_practice_at < ?", d.ID, now).
			Order("familiarity ASC").
			Limit(count).Find(&dungeonMonsters)
	}
	if err := query.Error; err != nil {
		return nil, err
	}

	log.Infof("got dungeon monsters %v", dungeonMonsters)
	return dungeonMonsters, nil
}

// GetDungeonMonsters - 获取当前 Dungeon 的 DungeonMonster，不会尝试解析 books 和 tags 的关联
func (d *Dungeon) GetDungeonMonsters(tx *gorm.DB, offset, limit int) ([]DungeonMonster, error) {
	var monsters []DungeonMonster
	err := tx.Where("dungeon_id = ?", d.ID).Order("item_id ASC").Offset(offset).Limit(limit).Find(&monsters).Error
	if err != nil {
		return nil, irr.Wrap(err, "failed to fetch item ids")
	}
	return monsters, nil
}

// GetDungeonMonstersWithExpandedAssociations - 获取当前 Dungeon 的 DungeonMonster 及其关联的 Items, Books, TagNames
func (d *Dungeon) GetDungeonMonstersWithExpandedAssociations(tx *gorm.DB, offset, limit int) ([]DungeonMonster, error) {
	// 获取关联的 Items, Books, TagNames
	bookIDs, items, tags, err := GetDungeonAssociations(tx, d.ID)
	if err != nil {
		return nil, err
	}

	// 直接获取 items
	itemSourceMap := make(map[utils.UInt64]utils.UInt64)
	for _, itemID := range items {
		itemSourceMap[itemID] = itemID
	}

	// 获取 book 关联的 items
	bookItemMap, err := GetItemIDsOfBooks(tx, bookIDs)
	if err != nil {
		return nil, err
	}
	for itemID := range bookItemMap {
		if bookID, exists := itemSourceMap[itemID]; exists {
			itemSourceMap[itemID] = bookID
		}
	}

	// 获取 tag 关联的 items
	tagItemMap, err := GetItemIDsOfTags(tx, tags)
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
	sort.Slice(itemIDs, func(i, j int) bool { // map 取值不稳定
		return itemIDs[i] < itemIDs[j]
	})

	// 获取所有 item 的详细信息并排序分页，不在内存里先裁剪的原因是如果查不到的话会导致列表 < limit
	// todo 当然还是有优化空间，比如空洞不多的情况下，先送内存裁剪的结果，有异常了再搜后续
	var itemsList []*Item
	if err = tx.Table("items").Where("id IN (?)", itemIDs).
		Offset(offset).Limit(limit).Find(&itemsList).Error; err != nil {
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
