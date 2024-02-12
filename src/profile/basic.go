// src/profile/profile_service.go
package profile

import (
	"github.com/bagaking/memorianexus/pkg/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SetProfileSettingsRequest 用于绑定请求正文
type SetProfileSettingsRequest struct {
	Nickname  string `json:"nickname,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

// SetProfileSettings 更新用户的昵称和头像URL
func (svr *Service) SetProfileSettings(c *gin.Context) {
	// 获取认证后的用户ID
	userID, exists := c.Get(auth.UserCtxKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// todo: validate at least one param exist
	var req SetProfileSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 更新用户配置信息
	if err := svr.Repo.UpdateProfileSettings(c, userID.(uint), req.Nickname, req.AvatarURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile settings updated successfully"})
}
