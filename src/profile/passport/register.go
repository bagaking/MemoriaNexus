package passport

import (
	"gorm.io/gorm"
)

// RegisterParams 定义注册参数结构体
type RegisterParams struct {
	Username string
	Email    string
	Password string
}

// RegisterService 提供用户注册的服务
type RegisterService struct {
	DB *gorm.DB
}

// NewRegisterService 创建一个新的RegisterService实例
func NewRegisterService(db *gorm.DB) *RegisterService {
	return &RegisterService{
		DB: db,
	}
}

// Register 创建一个新用户
func (service *RegisterService) Register(params RegisterParams) (*User, error) {
	// 通常你还需要在这里加密密码，这里为了简化示例，我们略过这一步
	// passwordHash := HashPassword(params.Password)

	user := &User{
		Username:     params.Username,
		Email:        params.Email,
		PasswordHash: params.Password, // 使用passwordHash代替明文密码
	}

	// 使用Gorm创建新的用户记录
	if err := service.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
