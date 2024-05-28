package dungeon

import (
	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/def"
)

type ReqCreateDungeon struct {
	DTODungeonFullData
}

type ReqUpdateDungeon struct {
	DTODungeonData
}

type RespUpdatedDungeon struct {
	DTODungeon
}

type DTODungeonFullData struct {
	DTODungeonData
	Books    []utils.UInt64 `json:"books,omitempty"`
	Items    []utils.UInt64 `json:"items,omitempty"`
	TagNames []string       `json:"tag_names,omitempty"`
	TagIDs   []utils.UInt64 `json:"tag_ids,omitempty"`
}

type DTODungeonData struct {
	Type        def.DungeonType `json:"type"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Rule        string          `json:"rule"`
}

type DTODungeon struct {
	ID utils.UInt64 `json:"id"`
	DTODungeonFullData
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type RespDungeonResults struct {
	TotalMonsters   int `json:"total_monsters"`
	TodayDifficulty int `json:"today_difficulty"`
}
