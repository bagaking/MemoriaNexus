package model

import (
	"errors"
	"time"

	"github.com/bagaking/memorianexus/src/def"

	"github.com/bagaking/goulp/wlog"
	"github.com/khicago/irr"

	"gorm.io/gorm"

	"github.com/bagaking/memorianexus/internal/utils"
)

type (
	Dungeon struct {
		ID     utils.UInt64 `gorm:"primaryKey;autoIncrement:false" json:"id"`
		UserID utils.UInt64 `gorm:"not null" json:"user_id,string"`

		Type        def.DungeonType `gorm:"not null" json:"type"`
		Title       string          `gorm:"not null" json:"title"`
		Description string          `json:"description"`
		Rule        string          `json:"rule"` // JSON format for detailed rule configuration

		CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
		UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

		DeletedAt gorm.DeletedAt
	}
)

type DungeonBook struct {
	DungeonID utils.UInt64 `gorm:"primaryKey"`
	BookID    utils.UInt64 `gorm:"primaryKey"`
}

type DungeonTag struct {
	DungeonID utils.UInt64 `gorm:"primaryKey"`
	TagID     utils.UInt64 `gorm:"primaryKey"`
}

// BeforeCreate 钩子
func (d *Dungeon) BeforeCreate(tx *gorm.DB) (err error) {
	// 确保UserID不为0
	if d.ID <= 0 {
		return errors.New("user UInt64 must be larger than zero")
	}
	return
}

// BeforeDelete hooks for cleaning up associations
func (d *Dungeon) BeforeDelete(tx *gorm.DB) (err error) {
	log := wlog.Common("BeforeDeleteDungeon")
	log.Infof("Deleting associations for dungeon ID %d", d.ID)

	if err = tx.Where("dungeon_id = ?", d.ID).Delete(&DungeonBook{}).Error; err != nil {
		return irr.Wrap(err, "failed to delete dungeon books")
	}

	if err = tx.Where("dungeon_id = ?", d.ID).Delete(&DungeonTag{}).Error; err != nil {
		return irr.Wrap(err, "failed to delete dungeon tags")
	}

	if err = tx.Where("dungeon_id = ?", d.ID).Delete(&DungeonMonster{}).Error; err != nil {
		return irr.Wrap(err, "failed to delete dungeon monsters")
	}

	return nil
}

func GetDungeonBookIDs(tx *gorm.DB, dungeonID utils.UInt64) ([]utils.UInt64, error) {
	var books []utils.UInt64

	tx = tx.Model(&DungeonBook{}).Select("book_id").Where("dungeon_id = ?", dungeonID)
	rows, err := tx.Rows()
	if err != nil {
		return nil, irr.Wrap(err, "failed to fetch book ids")
	}
	defer rows.Close()

	for rows.Next() {
		var id utils.UInt64
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		books = append(books, id)
	}
	return books, nil
}

func GetDungeonItemIDs(tx *gorm.DB, dungeonID utils.UInt64) ([]utils.UInt64, error) {
	var items []utils.UInt64

	tx = tx.Model(&DungeonMonster{}).Select("item_id").Where("dungeon_id = ?", dungeonID)
	rows, err := tx.Rows()
	if err != nil {
		return nil, irr.Wrap(err, "failed to fetch item ids")
	}
	defer rows.Close()

	for rows.Next() {
		var id utils.UInt64
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	return items, nil
}

func GetDungeonTagIDs(tx *gorm.DB, dungeonID utils.UInt64) ([]utils.UInt64, error) {
	var tags []utils.UInt64

	tx = tx.Model(&DungeonTag{}).Select("tag_id").Where("dungeon_id = ?", dungeonID)
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
		tags = append(tags, id)
	}
	return tags, nil
}

// GetDungeonAssociations Helper function to get associated books, items, and tags for a dungeon
func GetDungeonAssociations(tx *gorm.DB, dungeonID utils.UInt64) (books, items, tags []utils.UInt64, err error) {
	if books, err = GetDungeonBookIDs(tx, dungeonID); err != nil {
		return nil, nil, nil, irr.Wrap(err, "failed to fetch dungeon-book associations")
	}
	if tags, err = GetDungeonTagIDs(tx, dungeonID); err != nil {
		return nil, nil, nil, irr.Wrap(err, "failed to fetch dungeon-tag associations")
	}
	if items, err = GetDungeonItemIDs(tx, dungeonID); err != nil {
		return nil, nil, nil, irr.Wrap(err, "failed to fetch dungeon-item associations")
	}
	return books, items, tags, nil
}
