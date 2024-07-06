package item

import (
	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/def"
)

type (
	ReqCreateItem struct {
		Type       string              `json:"type"`
		Content    string              `json:"content"`
		Difficulty def.DifficultyLevel `json:"difficulty,omitempty"` // 难度，默认值为 NoviceNormal (0x01), todo: 考虑是否允许用户编辑，编辑后要引入写扩散
		Importance def.ImportanceLevel `json:"importance,omitempty"` // 重要程度，默认值为 DomainGeneral (0x01), todo: 考虑是否允许用户编辑
		BookIDs    []utils.UInt64      `json:"book_ids,omitempty"`   // 用于接收一个或多个 BookID
		Tags       []string            `json:"tags,omitempty"`       // 新增字段，用于接收一组 Tag 名称
	}

	ReqUpdateItem struct {
		Type       string              `json:"type,omitempty"`
		Content    string              `json:"content,omitempty"`
		Difficulty def.DifficultyLevel `json:"difficulty,omitempty"` // 难度，默认值为 NoviceNormal (0x01)
		Importance def.ImportanceLevel `json:"importance,omitempty"` // 重要程度，默认值为 DomainGeneral (0x01)
		Tags       []string            `json:"tags,omitempty"`       // 新增字段
	}

	ReqGetItems struct {
		UserID utils.UInt64 `form:"user_id" json:"user_id,omitempty"`
		BookID utils.UInt64 `form:"book_id" json:"book_id,omitempty"`
		Type   string       `form:"type" json:"type,omitempty"`
		Search string       `form:"search" json:"search,omitempty"`
	}

	ReqUploadItems struct {
		BookID *utils.UInt64 `form:"book_id,omitempty"`
	}
)

const (
	MaxBooksOncePerItem = 10 // 设定每个 Item 可以关联的最大 Books 数量
	MaxTagsOncePerItem  = 5  // 设定每个 Item 可以拥有的最大 Tags 数量
)
