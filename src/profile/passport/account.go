package passport

import (
	"gorm.io/gorm"
	"time"
)

// User 定义了用户账号管理模块的模型
type User struct {
	gorm.Model          // 包括了ID, CreatedAt, UpdatedAt, DeletedAt字段
	Username     string `gorm:"uniqueIndex;not null"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Provider     string `gorm:"default:local"` // 注册方式：local, google, facebook, twitter
	ProviderID   string // 社交账号的唯一标识
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
