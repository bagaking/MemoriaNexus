// File: src/profile/passport/registration_handler.go

package passport

import (
	"github.com/bagaking/memorianexus/pkg/auth"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RegisterRequest 定义注册请求的结构
type RegisterRequest struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

// ErrorMessage 定义错误信息的结构
type ErrorMessage struct {
	Message string `json:"message"`
}

// HandleRegister 处理注册请求
func (rs *RegisterService) HandleRegister(c *gin.Context) {
	var req RegisterRequest

	// 绑定请求体到结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// 验证两次输入的密码是否匹配
	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not encrypt password"})
		return
	}

	// 实例化注册服务并进行注册
	user, err := rs.Register(RegisterParams{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	})

	if err != nil {
		// 处理可能的数据库错误，如唯一性违反等
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	// todo: 这些值应该从配置中安全获取，现在 MVP 一下
	jwtService := auth.NewJWTService("my_secret_key", "MemoriaNexus")

	// 注册成功后生成JWT (short-ticket sample)
	token, err := jwtService.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// 设置Cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(72 * time.Hour),
		HttpOnly: true, // HttpOnly标志确保Javascript无法读取该cookie
	})

	// 返回创建成功的用户信息（注意不返回密码等敏感信息）
	c.JSON(http.StatusCreated, gin.H{"user": user, "token": token})
}
