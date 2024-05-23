package model

import (
	"time"

	"github.com/bagaking/goulp/wlog"
	"github.com/khicago/irr"
	"gorm.io/gorm"

	"github.com/bagaking/memorianexus/internal/utils"
)

type Item struct {
	ID     utils.UInt64 `gorm:"primaryKey;autoIncrement:true" json:"id"`
	UserID utils.UInt64 `gorm:"not null" json:"user_id,string"`

	Type    string
	Content string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ItemTag struct {
	ItemID utils.UInt64 `gorm:"primaryKey"`
	TagID  utils.UInt64 `gorm:"primaryKey"`
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
