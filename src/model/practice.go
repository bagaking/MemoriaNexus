package model

import (
	"context"
	"time"

	"golang.org/x/exp/rand"
	"gorm.io/gorm"

	"github.com/bagaking/goulp/wlog"
	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/internal/utils/cache"
	"github.com/bagaking/memorianexus/src/def"
	"github.com/khgame/memstore/cachekey"
)

const (
	balanceModeNewStuffRate utils.Percentage = 10
	defaultDalyNewCount                      = 10
	defaultThreshold                         = 3
)

type CParamNewStuffCountDaily struct {
	ID   utils.UInt64 `cachekey:"dungeon_id"`
	Date string       `cachekey:"date"`
}

var CKDungeonNewStuffCountDaily = cachekey.MustNewSchema[CParamNewStuffCountDaily](
	"dungeon:{dungeon_id}:new_stuff_count_daily:{date}", time.Hour*24) // 按天淘汰

func (d *Dungeon) GetMonstersForPractice(ctx context.Context, tx *gorm.DB, count int) ([]DungeonMonster, error) {
	log := wlog.ByCtx(ctx, "GetMonstersForPractice").
		WithField("dungeon_id", d.ID).
		WithField("count", count).
		WithField("quiz_mode", d.QuizMode).
		WithField("priority_mode", d.PriorityMode)

	log.Infof("start get monsters")
	now := time.Now()

	// Helper function to create the base query
	makeQuery := func(tx *gorm.DB, limit int) *gorm.DB {
		return tx.Where("dungeon_id = ? AND next_practice_at < ?", d.ID, now).
			Order("importance DESC, difficulty ASC").
			Limit(limit)
	}

	var dungeonMonsters []DungeonMonster
	var err error

	switch d.QuizMode {
	case def.QuizModeAlwaysNew:
		dungeonMonsters, err = getMonstersAlwaysNew(ctx, tx, makeQuery, count)
	case def.QuizModeAlwaysOld:
		dungeonMonsters, err = getMonstersAlwaysOld(ctx, tx, makeQuery, count)
	case def.QuizModeThreshold:
		dungeonMonsters, err = getMonstersThreshold(ctx, tx, makeQuery, d.ID, count)
	case def.QuizModeDynamic:
		dungeonMonsters, err = getMonstersDynamic(ctx, tx, makeQuery, d.ID, count)
	case def.QuizModeBalance:
		dungeonMonsters, err = getMonstersBalance(ctx, tx, makeQuery, count)
	default:
		dungeonMonsters, err = getMonstersBalance(ctx, tx, makeQuery, count)
	}

	if err != nil {
		return nil, err
	}

	log.Infof("got dungeon monsters %v", dungeonMonsters)
	return dungeonMonsters, nil
}

// Helper functions for different QuizModes

func getMonstersAlwaysNew(ctx context.Context, tx *gorm.DB, makeQuery func(*gorm.DB, int) *gorm.DB, count int) ([]DungeonMonster, error) {
	var m1, m2 []DungeonMonster
	if err := makeQuery(tx, count).Where("familiarity = 0").Find(&m1).Error; err != nil {
		return nil, err
	}
	if len(m1) < count {
		if err := makeQuery(tx, count-len(m1)).Where("familiarity > 0").Find(&m2).Error; err != nil {
			return nil, err
		}
	}
	return append(m1, m2...), nil
}

func getMonstersAlwaysOld(ctx context.Context, tx *gorm.DB, makeQuery func(*gorm.DB, int) *gorm.DB, count int) ([]DungeonMonster, error) {
	var m1, m2 []DungeonMonster
	if err := makeQuery(tx, count).Where("familiarity > 0").Find(&m1).Error; err != nil {
		return nil, err
	}
	if len(m1) < count {
		if err := makeQuery(tx, count-len(m1)).Where("familiarity = 0").Find(&m2).Error; err != nil {
			return nil, err
		}
	}
	return append(m1, m2...), nil
}

func getMonstersBalance(ctx context.Context, tx *gorm.DB, makeQuery func(*gorm.DB, int) *gorm.DB, count int) ([]DungeonMonster, error) {
	var m1, m2 []DungeonMonster
	if rand.Intn(100) < int(balanceModeNewStuffRate.Clamp0100()) {
		if err := makeQuery(tx, count).Where("familiarity = 0").Find(&m1).Error; err != nil {
			return nil, err
		}
	}
	if len(m1) < count {
		if err := makeQuery(tx, count-len(m1)).Where("familiarity > 0").Find(&m2).Error; err != nil {
			return nil, err
		}
	}
	return append(m1, m2...), nil
}

func getMonstersThreshold(ctx context.Context, tx *gorm.DB, makeQuery func(*gorm.DB, int) *gorm.DB, dungeonID utils.UInt64, count int) ([]DungeonMonster, error) {
	var m1, m2 []DungeonMonster
	minOne := min(defaultThreshold, count)
	if err := makeQuery(tx, minOne).Where("familiarity = 0").Find(&m1).Error; err != nil {
		return nil, err
	}
	if len(m1) < count {
		if err := makeQuery(tx, count-len(m1)).Where("familiarity > 0").Find(&m2).Error; err != nil {
			return nil, err
		}
	}
	return append(m1, m2...), nil
}

func getMonstersDynamic(ctx context.Context, tx *gorm.DB, makeQuery func(*gorm.DB, int) *gorm.DB, dungeonID utils.UInt64, count int) ([]DungeonMonster, error) {
	var m1, m2 []DungeonMonster

	key := CKDungeonNewStuffCountDaily.MustBuild(CParamNewStuffCountDaily{ID: dungeonID, Date: time.Now().Format("2006-01-02")})
	cNewStuff, err := cache.Client().Get(ctx, key).Int()
	if err != nil {
		wlog.ByCtx(ctx, "GetMonstersForPractice").
			WithError(err).WithField("cache_key", key).Warnf("get new stuff count failed")
		cNewStuff = defaultDalyNewCount
	}

	if cNewStuff < defaultDalyNewCount {
		if err = makeQuery(tx, count).Where("familiarity = 0").Find(&m1).Error; err != nil {
			return nil, err
		}
		cache.Client().Set(ctx, key, cNewStuff+1, time.Hour*24) // todo: 先按次数 set cache 实现了, 预期是应该在 submit 的时候再 incr cache
	}

	if len(m1) < count {
		if err = makeQuery(tx, count-len(m1)).Where("familiarity > 0").Find(&m2).Error; err != nil {
			return nil, err
		}
	}
	m2 = append(m2, m1...)
	if len(m2) < count { // 可能是没有走 threshold 的逻辑, 所以这里再查一次
		if err = makeQuery(tx, count-len(m2)).Where("familiarity = 0").Find(&m1).Error; err != nil {
			return nil, err
		}
		m2 = append(m2, m1...)
	}
	// todo: anyway 现在是一个临时的逻辑
	return m2, nil
}
