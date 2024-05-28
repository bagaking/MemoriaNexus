package dto

import (
	"time"

	"github.com/bagaking/memorianexus/src/model"

	"github.com/bagaking/memorianexus/internal/utils"
)

type (
	Book struct {
		ID          utils.UInt64 `json:"id"`
		UserID      utils.UInt64 `json:"user_id"`
		Title       string       `json:"title"`
		Description string       `json:"description"`
		CreatedAt   time.Time    `json:"created_at"`
		UpdatedAt   time.Time    `json:"updated_at"`
	}

	RespBookGet    = RespSuccess[*Book]
	RespBookCreate = RespSuccess[*Book]
)

func (b *Book) FromModel(m *model.Book) *Book {
	b.ID = m.ID
	b.UserID = m.UserID
	b.Title = m.Title
	b.Description = m.Description
	b.CreatedAt = m.CreatedAt
	b.UpdatedAt = m.UpdatedAt
	return b
}