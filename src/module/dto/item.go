package dto

import (
	"time"

	"github.com/bagaking/memorianexus/internal/utils"

	"github.com/bagaking/memorianexus/src/def"

	"github.com/bagaking/memorianexus/src/model"
)

// Item 数据传输对象
type (
	Item struct {
		ID         utils.UInt64        `json:"id"`
		CreatorID  utils.UInt64        `json:"creator_id"`
		Type       string              `json:"type"`
		Content    string              `json:"content"`
		Tags       []string            `json:"tags,omitempty"`
		CreatedAt  time.Time           `json:"created_at"`
		UpdatedAt  time.Time           `json:"updated_at"`
		Difficulty def.DifficultyLevel `json:"difficulty"`
		Importance def.ImportanceLevel `json:"importance"`
	}

	RespItemGet    = RespSuccess[*Item]
	RespItemDelete = RespSuccess[*Item]
	RespItemCreate = RespSuccess[*Item]
	RespItemUpdate = RespSuccess[*Item]
	RespItemList   = RespSuccessPage[*Item]
)

func (dto *Item) FromModel(item *model.Item, tags ...string) *Item {
	if item == nil {
		return nil
	}
	dto.ID = item.ID
	dto.CreatorID = item.CreatorID
	dto.Type = item.Type
	dto.Content = item.Content
	dto.CreatedAt = item.CreatedAt
	dto.UpdatedAt = item.UpdatedAt
	dto.Difficulty = item.Difficulty
	dto.Importance = item.Importance

	if tags != nil && len(tags) > 0 {
		dto.Tags = tags
	}
	return dto
}
