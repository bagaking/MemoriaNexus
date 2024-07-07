package dto

import (
	"time"

	"github.com/bagaking/memorianexus/internal/utils"

	"github.com/bagaking/memorianexus/src/def"
	"github.com/bagaking/memorianexus/src/model"
)

type (
	DungeonData struct {
		Type        def.DungeonType `json:"type"`
		Title       string          `json:"title"`
		Description string          `json:"description"`

		*SettingsMemorization
	}

	Dungeon struct {
		ID utils.UInt64 `json:"id"`
		DungeonData
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`

		Books []utils.UInt64 `json:"books,omitempty"`
		Items []utils.UInt64 `json:"items,omitempty"`
		Tags  []string       `json:"tags,omitempty"`
	}

	DungeonMonster struct {
		DungeonID utils.UInt64 `json:"dungeon_id,omitempty"`
		ItemID    utils.UInt64 `json:"item_id,omitempty"`

		// 用于 runtime
		PracticeAt     time.Time `json:"practice_at,omitempty"`      // 上次复习时间的记录
		NextPracticeAt time.Time `json:"next_practice_at,omitempty"` // 下次复习时间
		PracticeCount  uint32    `json:"practice_count,omitempty"`   // 复习次数 (考虑到可能会有 merge 次数等逻辑，这里先用一个相对大的空间）

		// 以下为宽表内容，为了加速查询
		Familiarity utils.Percentage `json:"familiarity,omitempty"` // UserMonster 向 DungeonMonster 单项同步

		Difficulty def.DifficultyLevel `json:"difficulty"` // Item -> DungeonMonster 单向同步
		Importance def.ImportanceLevel `json:"importance"` // Item -> DungeonMonster 单向同步

		// 以下为游戏性相关内容，由 AI 生成
		Visibility  utils.Percentage `json:"visibility,omitempty"` // Visibility 显影程度，根据复习次数变化
		Avatar      string           `json:"avatar,omitempty"`     // 怪物头像
		Name        string           `json:"name,omitempty"`
		Description string           `json:"description,omitempty"`

		// system
		SourceType model.MonsterSource `json:"source_type,omitempty"` // 记录插入时来源，方便原路径修改删除等
		SourceID   utils.UInt64        `json:"source_id,omitempty"`   // 记录插入时来源，方便原路径修改删除等
		CreatedAt  time.Time           `json:"created_at"`
	}

	RespDungeonResults struct {
		TotalMonsters   int `json:"total_monsters"`
		TodayDifficulty int `json:"today_difficulty"`
	}

	RespDungeon     = RespSuccess[*Dungeon]
	RespDungeonList = RespSuccessPage[*Dungeon]

	RespMonsterUpdate = RespSuccess[Updater[*DungeonMonster]]
	RespMonsterGet    = RespSuccess[*DungeonMonster]
	RespMonsterList   = RespSuccessPage[*DungeonMonster]
)

func (dto *DungeonMonster) FromModel(dm model.DungeonMonster) *DungeonMonster {
	dto.DungeonID = dm.DungeonID
	dto.ItemID = dm.ItemID
	dto.SourceType = dm.SourceType
	dto.SourceID = dm.SourceID

	// 用于 runtime
	dto.PracticeAt = dm.PracticeAt
	dto.NextPracticeAt = dm.NextPracticeAt
	dto.PracticeCount = dm.PracticeCount

	// 以下为宽表内容，为了加速查询
	dto.Familiarity = dm.Familiarity
	dto.Difficulty = dm.Difficulty
	dto.Importance = dm.Importance
	dto.CreatedAt = dm.CreatedAt

	// 用于游戏性
	dto.Avatar = dm.Avatar
	dto.Name = dm.Name
	dto.Description = dm.Description

	return dto
}

func (d *Dungeon) FromModel(model *model.Dungeon) *Dungeon {
	d.ID = model.ID
	d.Type = model.Type
	d.Title = model.Title
	d.Description = model.Description
	memSetting := model.MemorizationSetting
	d.SettingsMemorization = &SettingsMemorization{
		ReviewInterval:       &memSetting.ReviewInterval,
		DifficultyPreference: &memSetting.DifficultyPreference,
		QuizMode:             &memSetting.QuizMode,
		PriorityMode:         &memSetting.PriorityMode,
	}
	d.CreatedAt = model.CreatedAt
	d.UpdatedAt = model.UpdatedAt
	return d
}
