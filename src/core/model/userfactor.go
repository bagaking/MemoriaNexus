package model

import "gorm.io/gorm"

// User 定义了用户账号管理模块的模型
type User struct {
	gorm.Model          // Includes ID, CreatedAt, UpdatedAt, DeletedAt fields
	Username     string `gorm:"uniqueIndex;not null"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	TwoFactorID  uint   `gorm:"column:two_factor_setting_id"` // Reference to two-factor settings
}
