package model

import (
	"errors"
	"time"

	"github.com/bagaking/memorianexus/internal/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/khicago/irr"
)

// Profile 定义了用户个人信息的模型
type Profile struct {
	ID utils.UInt64 `gorm:"primaryKey;autoIncrement:false"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Nickname  string `gorm:"nickname,size:255"`
	Email     string `gorm:"email,size:255;not null;unique"`
	AvatarURL string `gorm:"avatar_url,size:255"`
	Bio       string `gorm:"bio,type:text"`

	DeletedAt gorm.DeletedAt `gorm:"index"`

	settingsMemorization *ProfileMemorizationSetting
	settingsAdvance      *ProfileAdvanceSetting

	points *ProfilePoints
}

// ProfileMemorizationSetting 定义了用户记忆设置的模型
type ProfileMemorizationSetting struct {
	ID utils.UInt64 `gorm:"primaryKey;autoIncrement:false"`

	ReviewInterval       uint   `gorm:"type:int unsigned"`
	DifficultyPreference uint8  `gorm:"type:tinyint unsigned"`
	QuizMode             string `gorm:"size:255"`
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
func (p *Profile) UpdateProfile(db *gorm.DB, updateData *Profile) error {
	// 使用Clauses提供onConflict来避免select的更新为空值的字段
	return db.Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{"nickname", "email", "avatar_url", "bio"}),
	}).Model(p).Where("id = ?", p.ID).Updates(updateData).Error
}

// EnsureLoadProfile 从数据库中加载用户个人信息
func EnsureLoadProfile(db *gorm.DB, uid utils.UInt64) (*Profile, error) {
	p := &Profile{
		ID: uid,
	}
	result := db.FirstOrInit(p, p)
	if result.Error != nil {
		return nil, irr.Wrap(result.Error, "search for profile failed")
	}

	// 如果是新初始化的对象，保存到数据库
	if result.RowsAffected == 0 {
		if err := db.Save(p).Error; err != nil {
			return nil, irr.Wrap(err, "create profile failed")
		}
	}

	return p, nil
}

// EnsureLoadProfileSettingsMemorization 从数据库中"懒加载"用户记忆设置
func (p *Profile) EnsureLoadProfileSettingsMemorization(db *gorm.DB) (*ProfileMemorizationSetting, error) {
	if p.settingsMemorization != nil {
		return p.settingsMemorization, nil
	}

	p.settingsMemorization = &ProfileMemorizationSetting{ID: p.ID}
	result := db.FirstOrInit(p.settingsMemorization, ProfileMemorizationSetting{ID: p.ID})
	if result.Error != nil {
		return nil, result.Error
	}

	// 如果是新初始化的对象，保存到数据库
	if result.RowsAffected == 0 {
		if err := db.Save(p.settingsMemorization).Error; err != nil {
			return nil, err
		}
	}

	return p.settingsMemorization, nil
}

// EnsureLoadProfileSettingsAdvance 从数据库中"懒加载"用户高级设置
func (p *Profile) EnsureLoadProfileSettingsAdvance(db *gorm.DB) (*ProfileAdvanceSetting, error) {
	if p.settingsAdvance != nil {
		return p.settingsAdvance, nil
	}

	p.settingsAdvance = &ProfileAdvanceSetting{ID: p.ID}
	result := db.FirstOrInit(p.settingsAdvance, ProfileAdvanceSetting{ID: p.ID})
	if result.Error != nil {
		return nil, result.Error
	}

	// 如果是新初始化的对象，保存到数据库
	if result.RowsAffected == 0 {
		if err := db.Save(p.settingsAdvance).Error; err != nil {
			return nil, err
		}
	}

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
