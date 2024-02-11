package auth

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// jwtCustomClaims 包含JWT的声明
type jwtCustomClaims struct {
	UserID uint `json:"userId"`
	jwt.StandardClaims
}

// JWTService 提供JWT令牌的服务
type JWTService struct {
	secretKey string
	issuer    string
}

// NewJWTService 创建JWT服务的新实例
func NewJWTService(secretKey, issuer string) *JWTService {
	return &JWTService{
		secretKey: secretKey,
		issuer:    issuer,
	}
}

// GenerateToken 生成JWT令牌
func (s *JWTService) GenerateToken(userID uint) (string, error) {
	// 设置JWT声明
	claims := &jwtCustomClaims{
		userID, // 用户ID从数据库用户模型中带入
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(), // 举例：让令牌在72小时后过期
			Issuer:    s.issuer,
		},
	}

	// 使用HMAC SHA256算法进行令牌签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return t, nil
}

// ValidateToken 验证JWT令牌
func (s *JWTService) ValidateToken(tokenString string) (*jwt.Token, error) {
	// 解析JWT令牌
	token, err := jwt.ParseWithClaims(tokenString, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 在这里验证token方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secretKey), nil
	})

	// 可能存在解析错误或令牌无效错误
	if err != nil {
		return nil, err
	}

	// 返回验证通过的令牌
	if _, ok := token.Claims.(*jwtCustomClaims); ok && token.Valid {
		return token, nil
	} else {
		return nil, errors.New("invalid token")
	}
}
