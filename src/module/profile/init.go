package profile

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
	// Other service dependencies would be added here
}

// NewService creates a new service instance with dependencies wired in.
func NewService(repo *gorm.DB) *Service {
	return &Service{db: repo}
}

// ApplyMux registers the service endpoints with the router.
// Swagger annotations describe the endpoints for the swagger documentation.
func (svr *Service) ApplyMux(router gin.IRouter) {
	router.GET("/me", svr.GetUserProfile)
	router.PUT("/me", svr.UpdateUserProfile)

	router.GET("/settings/memorization", svr.GetUserSettingsMemorization)
	router.PUT("/settings/memorization", svr.UpdateUserSettingsMemorization)
	router.GET("/settings/advance", svr.GetUserSettingsAdvance)
	router.PUT("/settings/advance", svr.GetUserSettingsAdvance)

	// todo: 严格来说这个不算是用户信息的部分。可以考虑分成系统级别的积分和游戏内的 -- 比如 get user points 时要触发挂机计算逻辑，应该只涉及游戏内的。
	router.GET("/points", svr.GetUserPoints)
}
