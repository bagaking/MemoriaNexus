package model

import (
	"context"
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

	if err = tx.Where("book_id = ?", b.ID).Delete(&BookItem{}).Error; err != nil {
		return irr.Wrap(err, "failed to delete book items")
	}

	return nil
}

func FindBook(ctx context.Context, tx *gorm.DB, id utils.UInt64) (*Book, error) {
	book := &Book{}
	result := tx.Where("id = ?", id).First(book)
	if err := result.Error; err != nil {
		return nil, err
	}
	return book, nil
}

func (b *Book) GetTagsName(ctx context.Context, tx *gorm.DB) ([]string, error) {
	return GetTagsByEntity(ctx, tx, b.ID)
}

func (b *Book) MPutItems(ctx context.Context, tx *gorm.DB, itemIDs []utils.UInt64) (successItemIDs []utils.UInt64, err error) {
	successItemIDs = make([]utils.UInt64, 0, len(itemIDs))
	for _, id := range itemIDs {
		bookItem := &BookItem{
			BookID: b.ID,
			ItemID: id,
		}
		if err = tx.Where(bookItem).FirstOrCreate(bookItem).Error; err != nil {
			return successItemIDs, irr.Wrap(err, "failed to add item to book")
		}
		successItemIDs = append(successItemIDs, id)
	}
	return successItemIDs, nil
}
