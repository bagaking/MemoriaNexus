package dto

import (
	"time"

	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/def"
	"github.com/bagaking/memorianexus/src/model"
)

type (
	DungeonFullData struct {
		DungeonData
		Books    []utils.UInt64 `json:"books,omitempty"`
		Items    []utils.UInt64 `json:"items,omitempty"`
		TagNames []string       `json:"tag_names,omitempty"`
		TagIDs   []utils.UInt64 `json:"tag_ids,omitempty"`
	}

	DungeonData struct {
		Type        def.DungeonType `json:"type"`
		Title       string          `json:"title"`
		Description string          `json:"description"`
		Rule        string          `json:"rule"`
	}

	Dungeon struct {
		ID utils.UInt64 `json:"id"`
		DungeonFullData
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	DungeonMonster struct {
		*Item

		Visibility int // 显影程度，根据复习次数变化

		DungeonID             utils.UInt64
		MonsterSource         model.MonsterSource
		MonsterSourceEntityID utils.UInt64
	}

	RespDungeonResults struct {
		TotalMonsters   int `json:"total_monsters"`
		TodayDifficulty int `json:"today_difficulty"`
	}

	RespDungeon     = RespSuccess[*Dungeon]
	RespDungeonList = RespSuccessPage[*Dungeon]

	RespMonsterGet  = RespSuccess[*DungeonMonster]
	RespMonsterList = RespSuccessPage[*DungeonMonster]
)

func (dto *DungeonMonster) FromModel(dm model.DungeonMonster, tags ...string) *DungeonMonster {
	dto.MonsterSource = dm.SourceType
	dto.MonsterSourceEntityID = dm.SourceID
	dto.DungeonID = dm.DungeonID

	dto.Item = dto.Item.FromModel(dm.Monster(), tags...)

	return dto
}

func (d *Dungeon) FromModel(model *model.Dungeon) *Dungeon {
	d.ID = model.ID
	d.Type = model.Type
	d.Title = model.Title
	d.Description = model.Description
	d.Rule = model.Rule
	d.CreatedAt = model.CreatedAt
	d.UpdatedAt = model.UpdatedAt
	return d
}
