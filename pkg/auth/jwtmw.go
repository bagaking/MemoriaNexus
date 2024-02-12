// File: src/app/gw/middleware.go

package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	BearerSchema        = "Bearer "
	HeaderAuthorization = "Authorization"
	UserCtxKey          = "UserID"
)

// AuthMiddleware 是一个检查JWT是否有效的中间件
func (s *JWTService) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(HeaderAuthorization)
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		value := strings.TrimPrefix(header, BearerSchema)
		if value == header {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization schema is wrong"})
			return
		}

		token, err := s.ValidateToken(value)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
			c.Set(UserCtxKey, claims.UserID)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		}
	}
}
