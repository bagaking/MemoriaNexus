package model

import (
	"gorm.io/gorm"
)

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

// RegisterParams 定义注册参数结构体
type RegisterParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register 创建一个新用户
func (repo *Repo) Register(params RegisterParams) (*User, error) {
	// 通常你还需要在这里加密密码，这里为了简化示例，我们略过这一步
	// passwordHash := HashPassword(params.Password)

	user := &User{
		Username:     params.Username,
		Email:        params.Email,
		PasswordHash: params.Password, // 使用passwordHash代替明文密码
	}

	// 使用Gorm创建新的用户记录
	if err := repo.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// FindUserByName 用过名字找到一个新用户
func (repo *Repo) FindUserByName(username string) (*User, error) {
	var user User
	err := repo.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
