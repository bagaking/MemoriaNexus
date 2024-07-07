package model

import (
	"context"
	"errors"
	"time"

	"github.com/bagaking/memorianexus/internal/utils"

	"github.com/bagaking/memorianexus/src/def"

	"github.com/bagaking/goulp/wlog"
	"github.com/khicago/irr"

	"gorm.io/gorm"
)

type (
	Dungeon struct {
		ID     utils.UInt64 `gorm:"primaryKey;autoIncrement:false" json:"id"`
		UserID utils.UInt64 `gorm:"not null" json:"user_id,string"`

		Type        def.DungeonType `gorm:"not null" json:"type"`
		Title       string          `gorm:"not null" json:"title"`
		Description string          `json:"description"`

		MemorizationSetting

		// system
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

func CreateDungeon(ctx context.Context, tx *gorm.DB, d *Dungeon) (*Dungeon, error) {
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()
	if err := tx.Create(d).Error; err != nil {
		return nil, irr.Wrap(err, "failed to create dungeon")
	}
	return d, nil
}

func FindDungeon(ctx context.Context, tx *gorm.DB, dungeonID utils.UInt64) (*Dungeon, error) {
	dungeon := &Dungeon{}
	if err := tx.Where("id = ?", dungeonID).First(dungeon).Error; err != nil {
		return nil, err
	}
	return dungeon, nil
}

func (d *Dungeon) GetBookIDs(ctx context.Context, tx *gorm.DB) ([]utils.UInt64, error) {
	var books []utils.UInt64

	tx = tx.Model(&DungeonBook{}).Select("book_id").Where("dungeon_id = ?", d.ID)
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

func (d *Dungeon) GetItemIDs(ctx context.Context, tx *gorm.DB) ([]utils.UInt64, error) {
	var items []utils.UInt64
	tx = tx.Model(&DungeonMonster{}).Select("item_id").Where("dungeon_id = ?", d.ID)
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

func (d *Dungeon) GetTags(ctx context.Context, tx *gorm.DB) ([]string, error) {
	tags, err := GetTagsByEntity(ctx, tx, d.ID)
	if err != nil {
		return nil, irr.Wrap(err, "failed to fetch tag ids")
	}
	return tags, nil
}

// GetAssociations Helper function to get associated books, items, and tags for a dungeon
func (d *Dungeon) GetAssociations(ctx context.Context, tx *gorm.DB) (books, items []utils.UInt64, tags []string, err error) {
	if books, err = d.GetBookIDs(ctx, tx); err != nil {
		return nil, nil, nil, irr.Wrap(err, "failed to fetch dungeon-book associations")
	}
	if tags, err = d.GetTags(ctx, tx); err != nil {
		return nil, nil, nil, irr.Wrap(err, "failed to fetch dungeon-tag associations")
	}
	if items, err = d.GetItemIDs(ctx, tx); err != nil { // todo: 先不分页
		return nil, nil, nil, irr.Wrap(err, "failed to fetch dungeon-item associations")
	}
	return books, items, tags, nil
}

func (d *Dungeon) SubtractBooks(ctx context.Context, tx *gorm.DB, books []utils.UInt64) (successIDs []utils.UInt64, err error) {
	successIDs = make([]utils.UInt64, 0, len(books))
	for _, bookID := range books {
		// 删除关联
		if err = tx.Where("dungeon_id = ? AND book_id = ?", d.ID, bookID).Delete(&DungeonBook{}).Error; err != nil {
			return successIDs, err
		}
		successIDs = append(successIDs, bookID)
	}
	return successIDs, nil
}
