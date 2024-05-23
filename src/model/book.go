package model

import (
	"time"

	"github.com/bagaking/goulp/wlog"
	"github.com/khicago/irr"

	"github.com/bagaking/memorianexus/internal/utils"

	"gorm.io/gorm"
)

type Book struct {
	ID          utils.UInt64 `gorm:"primaryKey;autoIncrement:false" json:"id"`
	UserID      utils.UInt64 `gorm:"index:idx_user;not null" json:"user_id"`
	Title       string       `gorm:"not null"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`

	Tags  []*Tag  `gorm:"many2many:BookTag;"`
	Items []*Item `gorm:"many2many:BookItem;"`
}

// BeforeDelete is a GORM hook that is called before deleting a book.
func (b *Book) BeforeDelete(tx *gorm.DB) (err error) {
	log := wlog.Common("BeforeDeleteBook")
	log.Infof("Deleting associations for book ID %d", b.ID)

	// todo: refine this

	if err = tx.Where("book_id = ?", b.ID).Delete(&BookTag{}).Error; err != nil {
		return irr.Wrap(err, "failed to delete book tags")
	}

	if err = tx.Where("book_id = ?", b.ID).Delete(&BookItem{}).Error; err != nil {
		return irr.Wrap(err, "failed to delete book items")
	}

	return nil
}

type BookTag struct {
	BookID utils.UInt64 `gorm:"primaryKey"`
	TagID  utils.UInt64 `gorm:"primaryKey"`
}

type BookItem struct {
	BookID utils.UInt64 `gorm:"primaryKey"`
	ItemID utils.UInt64 `gorm:"primaryKey"`
}
