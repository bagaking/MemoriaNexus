package def

import (
	"database/sql/driver"
	"errors"

	"github.com/bytedance/sonic"
)

type (
	// PriorityMode - 复习时的出场顺序
	PriorityMode []PriorityModeSetting

	// PriorityModeSetting - 复习时的出场顺序
	PriorityModeSetting string
)

const (
	// PriorityModeFamiliarityASC - 熟练度: 生疏优先
	PriorityModeFamiliarityASC PriorityModeSetting = "familiarity_asc"
	// PriorityModeFamiliarityDESC - 熟练度: 熟悉优先
	PriorityModeFamiliarityDESC PriorityModeSetting = "familiarity_desc"

	// PriorityModeTimePassDESC - 近远期: 近期优先
	PriorityModeTimePassDESC PriorityModeSetting = "time_pass_desc"
	// PriorityModeTimePassASC - 近远期: 远期优先
	PriorityModeTimePassASC PriorityModeSetting = "time_pass_asc"

	// PriorityModeDifficultyASC - 难度: 简单优先
	PriorityModeDifficultyASC PriorityModeSetting = "difficulty_asc"
	// PriorityModeDifficultyDESC - 难度: 困难优先
	PriorityModeDifficultyDESC PriorityModeSetting = "difficulty_desc"

	// PriorityModeRelatedASC - 关联度: 关联度高的优先
	PriorityModeRelatedASC PriorityModeSetting = "related_asc"
	// PriorityModeImportanceASC - 重要度: 重要程度高的优先
	PriorityModeImportanceASC PriorityModeSetting = "importance_asc"
	// PriorityModeShuffle - 是否打乱, 不打乱就默认按 DB 顺序
	PriorityModeShuffle PriorityModeSetting = "shuffle"
)

// Scan 实现 sql.Scanner 接口
func (r *PriorityMode) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return sonic.Unmarshal(bytes, r)
}

// Value 实现 driver.Valuer 接口
func (r PriorityMode) Value() (driver.Value, error) {
	return sonic.Marshal(r)
}
