package model

import (
	"time"

	"github.com/bagaking/memorianexus/internal/util"
)

type Item struct {
	ID     util.UInt64 `gorm:"primaryKey;autoIncrement:true" json:"id"`
	UserID util.UInt64 `gorm:"not null" json:"user_id,string"`

	Type    string
	Content string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ItemTag struct {
	ItemID util.UInt64 `gorm:"primaryKey"`
	TagID  util.UInt64 `gorm:"primaryKey"`
}
