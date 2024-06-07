package dungeon

import (
	"context"
	"time"

	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/def"
	"github.com/bagaking/memorianexus/src/model"
)

// calculateNextPracticeAt calculates the next recall time for a monster
func calculateNextPracticeAt(ctx context.Context,
	familiarity utils.Percentage,
	importance def.ImportanceLevel,
	difficulty def.DifficultyLevel,
	userSettings *model.ProfileMemorizationSetting,
	shouldPracticeAt time.Time,
) time.Time {
	setting := model.DefaultMemorizationSetting
	if userSettings != nil {
		setting = *userSettings
	}

	// 获取下次复习的时间间隔
	nextInterval := setting.ReviewIntervalSetting.GetInterval(familiarity)

	// 根据重要性调整间隔时间，重要性越高，间隔时间越短
	importanceFactor := 1 / (1 + importance.Normalize())

	// 根据难度调整间隔时间，难度越高，间隔时间越短
	difficultyFactor := 1 / (1 + difficulty.Normalize())

	// 根据用户设置调整复习间隔
	nextInterval = time.Duration(float64(nextInterval) * importanceFactor * difficultyFactor)

	// 根据用户的难度偏好调整复习间隔
	difficultyPreferenceFactor := 1 + float64(setting.DifficultyPreference)/10
	nextInterval = time.Duration(float64(nextInterval) * difficultyPreferenceFactor)

	// 考虑用户复习延迟的情况
	now := time.Now()
	if shouldPracticeAt.Add(nextInterval).Before(now) {
		// 如果用户延迟复习，超过当前应该复习的档位 （也就是如果按时复习，现在应该复习超过一次了），则缩短下次复习时间
		nextInterval = nextInterval / 2
		// 由于创建 DungeonMonster 时不一定复习，因此可能导致第一次时间普遍缩短，先认为可以接受
	}

	return now.Add(nextInterval)
}
