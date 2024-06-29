package dto

import (
	"time"

	"github.com/bagaking/memorianexus/internal/utils"

	"github.com/bagaking/memorianexus/src/model"
)

// models
type (
	// Book - dto model for book
	Book struct {
		ID          utils.UInt64 `json:"id"`
		UserID      utils.UInt64 `json:"user_id"`
		Title       string       `json:"title"`
		Description string       `json:"description"`
		CreatedAt   time.Time    `json:"created_at"`
		UpdatedAt   time.Time    `json:"updated_at"`

		Tags []string `json:"tags,omitempty"`
	}

	RespBookList = RespSuccessPage[*Book]
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

func (b *Book) SetTags(tags []string) *Book {
	b.Tags = tags
	return b
}

// responses
type (
	RespBookGet    = RespSuccess[*Book]
	RespBookCreate = RespSuccess[*Book]
	RespBookUpdate = RespSuccess[*Book]

	RespBookDelete = RespSuccess[utils.UInt64]
	RespBooks      = RespSuccessPage[Book]
)
