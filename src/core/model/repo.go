package model

import "gorm.io/gorm"

// Repo 提供用户注册的服务
type Repo struct {
	DB *gorm.DB
}

// NewRepo 创建一个新的RegisterService实例
func NewRepo(db *gorm.DB) *Repo {
	return &Repo{
		DB: db,
	}
}
