package item

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

func (svr *Service) ApplyMux(group gin.IRouter) {
	group.POST("", svr.CreateItem)
	group.GET("", svr.GetItems)
	group.GET("/:id", svr.GetItem)
	group.PUT("/:id", svr.UpdateItem)
	group.DELETE("/:id", svr.DeleteItem)
}
