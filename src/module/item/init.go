package item

import (
	"github.com/bagaking/memorianexus/internal/utils"
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
	group.POST("upload", svr.UploadItems)
	group.GET("", svr.GetItems)

	idGroup := group.Group("/:id").Use(utils.GinMWParseID())
	{
		idGroup.GET("", svr.ReadItem)
		idGroup.PUT("", svr.UpdateItem)
		idGroup.DELETE("", svr.DeleteItem)
	}
}
