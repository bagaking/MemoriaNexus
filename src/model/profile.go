package model

import (
	"context"
	"errors"
	"time"

	"github.com/bagaking/memorianexus/src/def"

	"github.com/bagaking/memorianexus/internal/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/khicago/irr"
)

// Profile 定义了用户个人信息的模型
type Profile struct {
	ID utils.UInt64 `gorm:"primaryKey;autoIncrement:false"`

	Nickname  string `gorm:"nickname,size:255"`
	Email     string `gorm:"email,size:255;not null;unique"`
	AvatarURL string `gorm:"avatar_url,size:255"`
	Bio       string `gorm:"bio,type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time

	DeletedAt gorm.DeletedAt `gorm:"index"`

	settingsMemorization *ProfileMemorizationSetting
	settingsAdvance      *ProfileAdvanceSetting

	points *ProfilePoints
}

// ProfileMemorizationSetting 定义了用户记忆设置的模型
type (

	// ProfileMemorizationSetting
	//
	// 在 User 的 ProfileMemorizationSetting 中，可以设置默认的 QuizMode 和 Priority
	// Dungeon 创建时，会从 User 的 ProfileMemorizationSetting 中获取默认的 QuizMode 和 Priority
	// 并独立于 User 的设置进行修改，后续 User 的设置不会覆盖 Dungeon 的设置
	// 可以在 Dungeon 的 MemorizationSetting 中主动选择从 User 的 ProfileMemorizationSetting 中同步
	ProfileMemorizationSetting struct {
		ID utils.UInt64 `gorm:"primaryKey;autoIncrement:false"`
		MemorizationSetting
	}

	MemorizationSetting struct {
		// 复习时间的配置，是一组时间，作为根据复习结算时的熟练度来选择下次复习时间的依据
		ReviewInterval def.RecallIntervalLevel `gorm:"type:string"`

		// 用户挑战偏好
		DifficultyPreference utils.Percentage `gorm:"type:tinyint unsigned"`

		// 倾向的战斗模式，决定了已经在时间内 monster 出场时，新增和复习的出现策略
		QuizMode def.QuizMode `gorm:"size:255"`

		// 倾向的战斗模式，决定了已经在时间内 monster 出场时，进行选择的优先级顺序
		PriorityMode def.PriorityMode `gorm:"size:255"`
	}
)

var DefaultMemorizationSetting = MemorizationSetting{
	// 用户的设置值
	ReviewInterval:       def.DefaultRecallIntervals, // 先用 day
	DifficultyPreference: 1,
	QuizMode:             def.QuizModeBalance,
	PriorityMode: def.PriorityMode{
		def.PriorityModeFamiliarityDESC,
		def.PriorityModeTimePassASC,
		def.PriorityModeDifficultyASC,
		def.PriorityModeImportanceASC,
	},
}

// ProfileAdvanceSetting 定义了用户高级设置的模型
type ProfileAdvanceSetting struct {
	ID utils.UInt64 `gorm:"primaryKey;autoIncrement:false"`

	Theme              string `gorm:"theme,size:255;default:'light'"`
	Language           string `gorm:"language,size:255;default:'en'"`
	EmailNotifications bool   `gorm:"email_notifications,default:true"`
	PushNotifications  bool   `gorm:"push_notifications,default:true"`
}

// BeforeCreate 钩子
func (p *Profile) BeforeCreate(tx *gorm.DB) (err error) {
	// 确保UserID不为0
	if p.ID <= 0 {
		return errors.New("user UInt64 must be larger than zero")
	}
	return
}

// BeforeUpdate 钩子
func (p *Profile) BeforeUpdate(tx *gorm.DB) (err error) {
	// 更新逻辑（如果需要的话）
	return
}

// GetSettingsMemorization 检索用户记忆设置
func (p *Profile) GetSettingsMemorization() *ProfileMemorizationSetting {
	return p.settingsMemorization
}

// GetSettingsAdvance 检索用户高级设置
func (p *Profile) GetSettingsAdvance() *ProfileAdvanceSetting {
	return p.settingsAdvance
}

// GetPoints 检索用户点数
func (p *Profile) GetPoints() *ProfilePoints {
	return p.points
}

// CreateProfile 创建新的用户资料
func (p *Profile) CreateProfile(db *gorm.DB) error {
	return db.Create(p).Error
}

// UpdateProfile 更新现有用户资料
func (p *Profile) UpdateProfile(db *gorm.DB) error {
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}}, // 指定冲突的列
		UpdateAll: true,                          // 在冲突时更新所有列
	}).Create(p).Error
}

// FindProfile 读取个人信息
func FindProfile(ctx context.Context, db *gorm.DB, uid utils.UInt64) (*Profile, error) {
	cond := &Profile{ID: uid}
	result := db.Where(cond).First(cond)
	if result.Error != nil {
		return nil, result.Error
	}
	return cond, nil
}

// EnsureProfile 从数据库中加载用户个人信息
func EnsureProfile(ctx context.Context, db *gorm.DB, uid utils.UInt64) (*Profile, error) {
	cond := &Profile{
		ID:        uid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	result := db.Where("id = ?", cond.ID).FirstOrCreate(cond)
	if result.Error != nil {
		return nil, irr.Wrap(result.Error, "search for profile failed")
	}
	return cond, nil
}

func (p *Profile) GetSettingsMemorizationOrDefault(ctx context.Context, tx *gorm.DB) (*ProfileMemorizationSetting, error) {
	var userSettings ProfileMemorizationSetting
	if err := tx.Where("id = ?", p.ID).First(&userSettings).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, irr.Wrap(err, "failed to fetch user settings, user_id=%v", p.ID)
		}
		userSettings.MemorizationSetting = DefaultMemorizationSetting
		userSettings.ID = p.ID
	} else {
		p.settingsMemorization = &userSettings
	}
	return &userSettings, nil
}

// EnsureSettingsMemorization 从数据库中"懒加载"用户记忆设置
func (p *Profile) EnsureSettingsMemorization(db *gorm.DB) (*ProfileMemorizationSetting, error) {
	if p.settingsMemorization != nil {
		return p.settingsMemorization, nil
	}

	cond := &ProfileMemorizationSetting{
		ID:                  p.ID,
		MemorizationSetting: DefaultMemorizationSetting,
	}
	result := db.Where("id = ?", p.ID).FirstOrCreate(cond)
	if result.Error != nil {
		return nil, result.Error
	}

	p.settingsMemorization = cond
	return p.settingsMemorization, nil
}

// EnsureLoadProfileSettingsAdvance 从数据库中"懒加载"用户高级设置
func (p *Profile) EnsureLoadProfileSettingsAdvance(db *gorm.DB) (*ProfileAdvanceSetting, error) {
	if p.settingsAdvance != nil {
		return p.settingsAdvance, nil
	}

	cond := &ProfileAdvanceSetting{ID: p.ID}
	result := db.Where("id = ?", cond.ID).FirstOrCreate(cond)
	if result.Error != nil {
		return nil, result.Error
	}
	p.settingsAdvance = cond
	return p.settingsAdvance, nil
}

// UpdateSettingsMemorization 使用保存逻辑更新用户记忆设置
func (p *Profile) UpdateSettingsMemorization(db *gorm.DB, updater *ProfileMemorizationSetting) error {
	if updater == nil {
		return irr.Trace("updater cannot be nil")
	}
	if updater.ID != p.ID {
		// 确保settingsMemorization的ID与Profile的ID匹配
		return errors.New("profile UInt64 does not match with settingsMemorization UInt64")
	}

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(updater)

	return result.Error
}

// UpdateSettingsAdvance 使用保存逻辑更新用户高级设置
func (p *Profile) UpdateSettingsAdvance(db *gorm.DB, updater *ProfileAdvanceSetting) error {
	if updater == nil {
		return irr.Trace("updater cannot be nil")
	}
	if updater.ID != p.ID {
		// 确保settingsAdvance的ID与Profile的ID匹配
		return irr.Trace("profile UInt64 does not match with settingsAdvance UInt64")
	}

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(updater)

	return result.Error
}
