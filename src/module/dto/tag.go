package dto

import (
	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/model"
)

type (
	Tag struct {
		ID   utils.UInt64 `json:"id"`
		Name string       `json:"name"`
	}

	RespTagGet  = RespSuccess[*Tag]
	RespTagList = RespSuccessPage[*Tag]
)

func (t *Tag) FromModel(tag *model.Tag) *Tag {
	t.ID = tag.ID
	t.Name = tag.Name
	return t
}
