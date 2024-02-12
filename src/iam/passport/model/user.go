package model

import (
	"gorm.io/gorm"
	"time"
)

// User 定义了用户账号管理模块的模型
type User struct {
	gorm.Model          // Includes ID, CreatedAt, UpdatedAt, DeletedAt fields
	Username     string `gorm:"uniqueIndex;not null"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	TwoFactorID  uint   `gorm:"column:two_factor_setting_id"` // Reference to two-factor settings
}

type OAuthCredential struct {
	gorm.Model
	UserID     uint   `gorm:"index;not null"`
	User       User   `gorm:"foreignKey:UserID;references:ID"`
	Provider   string `gorm:"not null"` // local, google, facebook, twitter, wechat
	ProviderID string `gorm:"index"`
}

type TwoFactorSetting struct {
	gorm.Model
	IsEnabled      bool
	Phone          string
	SecondaryEmail string
	BackupCodes    string // JSON array of backup codes
}

// BeforeCreate 是Gorm的hook，在创建记录前调用
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = time.Now()
	}
	return nil
}

// BeforeUpdate 是Gorm的hook，在更新记录前调用
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()
	return nil
}
