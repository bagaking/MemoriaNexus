package campaign

import (
	"context"
	"time"

	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/def"
	"github.com/bagaking/memorianexus/src/model"
)

// FamiliarityFixSettings represents user-specific settings for familiarity calculation
type (
	FamiliarityFixSettings struct {
		PastWeight    float64
		CurrentWeight float64
		DecaySetting  DecaySetting
	}

	DecaySetting map[int]float64
)

func (ds DecaySetting) Factor(hours float64) float64 {
	factor := 1.0
	for h, f := range DefaultDecaySettings.DecaySetting {
		if hours > float64(h) {
			factor = f
		}
	}
	return factor
}

var DefaultDecaySettings = FamiliarityFixSettings{
	PastWeight:    0.8,
	CurrentWeight: 0.2,
	DecaySetting: DecaySetting{
		24: 0.9,
		48: 0.8,
	},
}

// CalculateNewFamiliarity 计算熟练度
// 熟练度主要受到难度和复习间隔影响，难度影响新熟练度的权重，复习间隔影响旧熟练度的比重
// 熟练度的计算反应的是用户固有的记忆效果，和任务、偏好等外在因素无关
func CalculateNewFamiliarity(currentFamiliarity, damageRate utils.Percentage, lastPracticeAt time.Time, difficulty def.DifficultyLevel) utils.Percentage {
	// 时间衰减因子
	timeDecayFactor := DefaultDecaySettings.DecaySetting.Factor(time.Since(lastPracticeAt).Hours())
	// 难度因子
	difficultyFactor := difficulty.Factor()

	// 动态调整过往熟练度和当前熟练度的比例
	pastWeight := DefaultDecaySettings.PastWeight * timeDecayFactor        // Decay 修正旧值占比
	currentWeight := DefaultDecaySettings.CurrentWeight * difficultyFactor // 难度越高，新值影响越大

	// 归一化权重
	totalWeight := pastWeight + currentWeight
	pastWeight /= totalWeight
	currentWeight /= totalWeight

	// 设置最小权重阈值，防止旧值的占比过低
	const minPastWeight = 0.4
	if pastWeight < minPastWeight {
		pastWeight = minPastWeight
		currentWeight = 1.0 - pastWeight
	}

	// 计算新的熟练度
	newRate := currentFamiliarity.NormalizedFloat()*pastWeight + damageRate.NormalizedFloat()*currentWeight
	if newRate > 1.0 {
		newRate = 1.0
	}
	newFamiliarity := utils.Percentage(newRate * 100)

	return newFamiliarity
}

// CalculateNextPracticeAt calculates the next recall time for a monster
// 下次练习时间主要由 familiarity 决定，受到重要程度、用户挑战偏好和 Dungeon 挑战偏好影响而进行修正
// 下次练习时间计算是用户在固有的记忆效果的基础上，反应任务、偏好等外在因素的要求
func CalculateNextPracticeAt(ctx context.Context,
	familiarity utils.Percentage,
	importance def.ImportanceLevel,
	memSetting *model.MemorizationSetting,
) time.Time {
	setting := model.DefaultMemorizationSetting
	if memSetting != nil {
		setting = *memSetting
	}

	// 获取下次复习的时间间隔
	nextInterval := setting.ReviewInterval.GetInterval(familiarity)

	// 根据重要性调整间隔时间，重要性越高，间隔时间越短
	importanceFactor := 1 / (1 + importance.Normalize())
	nextInterval = time.Duration(float64(nextInterval) * importanceFactor) // 最多可能调整到约 3/4

	// 根据用户的难度偏好调整复习间隔
	difficultyPreferenceFactor := 1 / (1 + setting.DifficultyPreference.NormalizedFloat()) // 最多可能调整到越 1/2
	nextInterval = time.Duration(float64(nextInterval) * difficultyPreferenceFactor)

	if nextInterval < time.Minute*3 {
		nextInterval = time.Minute * 3 // 最少 3 分钟，之后可以改成配置
	}

	// 计算下次复习时间
	return time.Now().Add(nextInterval)
}
