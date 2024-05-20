package model

import (
	"time"

	"github.com/bagaking/memorianexus/internal/util"

	"gorm.io/gorm"
)

type Book struct {
	ID          util.UInt64 `gorm:"primaryKey;autoIncrement:false" json:"id"`
	UserID      util.UInt64 `gorm:"index:idx_user;not null" json:"user_id"`
	Title       string      `gorm:"not null"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	Tags  []*Tag  `gorm:"many2many:BookTag;"`
	Items []*Item `gorm:"many2many:BookItem;"`
}

type BookTag struct {
	BookID util.UInt64 `gorm:"primaryKey"`
	TagID  util.UInt64 `gorm:"primaryKey"`
}

type BookItem struct {
	BookID util.UInt64 `gorm:"primaryKey"`
	ItemID util.UInt64 `gorm:"primaryKey"`
}
