package dto

import (
	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/model"
)

// Monster 数据传输对象
type (
	Monster struct {
		*Item

		Visibility int // 显影程度，根据复习次数变化

		DungeonID             utils.UInt64
		MonsterSource         model.MonsterSource
		MonsterSourceEntityID utils.UInt64
	}

	RespMonsterGet  = RespSuccess[*Monster]
	RespMonsterList = RespSuccessPage[*Monster]
)

func (dto *Monster) FromModel(dm model.DungeonMonster, tags ...string) *Monster {
	dto.MonsterSource = dm.SourceType
	dto.MonsterSourceEntityID = dm.SourceID
	dto.DungeonID = dm.DungeonID

	dto.Item = dto.Item.FromModel(dm.Monster(), tags...)

	return dto
}
