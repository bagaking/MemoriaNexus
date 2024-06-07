package model

import (
	"errors"
	"time"

	"github.com/bagaking/memorianexus/internal/utils"

	"github.com/bagaking/goulp/wlog"
	"github.com/khicago/irr"

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

// BeforeCreate 钩子
func (b *Book) BeforeCreate(tx *gorm.DB) (err error) {
	// 确保UserID不为0
	if b.ID <= 0 {
		return errors.New("user UInt64 must be larger than zero")
	}
	return
}

func (b *Book) TableName() string {
	return "books"
}

type BookTag struct {
	BookID utils.UInt64 `gorm:"primaryKey"`
	TagID  utils.UInt64 `gorm:"primaryKey"`
}

func (BookTag) Associate(bookID, tagID utils.UInt64) ITagAssociate {
	return BookTag{BookID: bookID, TagID: tagID}
}

func (BookTag) Type() TagRefType {
	return BookTagRef
}

var _ ITagAssociate = &BookTag{}

func (b BookTag) TableName() string {
	return "book_tags"
}

type BookItem struct {
	BookID utils.UInt64 `gorm:"primaryKey"`
	ItemID utils.UInt64 `gorm:"primaryKey"`
}

func (b *BookItem) TableName() string {
	return "book_items"
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
