package model

import (
	"gorm.io/gorm"
)

// Profile 定义了用户配置的模型
type Profile struct {
	gorm.Model
	UserID    uint   `gorm:"uniqueIndex"` // 用户ID
	Nickname  string `gorm:"size:255"`    // 用户昵称
	AvatarURL string `gorm:"size:255"`    // 用户头像URL
}

// BeforeCreate 钩子
func (p *Profile) BeforeCreate(tx *gorm.DB) (err error) {
	// 插入逻辑（如果需要的话）
	return
}

// BeforeUpdate 钩子
func (p *Profile) BeforeUpdate(tx *gorm.DB) (err error) {
	// 更新逻辑（如果需要的话）
	return
}
