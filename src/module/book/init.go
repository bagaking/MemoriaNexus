package book

import (
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

	"github.com/bagaking/memorianexus/internal/utils"
)

type Service struct {
	db *gorm.DB
}

var svr *Service

func Init(db *gorm.DB) (*Service, error) {
	svr = &Service{
		db: db,
	}
	return svr, nil
}

func (svr *Service) ApplyMux(group gin.IRouter) {
	group.POST("", svr.CreateBook)
	group.GET("", svr.ListBooks)
	idGroup := group.Group("/:id").Use(utils.GinMWParseID())
	{
		idGroup.GET("", svr.ReadBook)
		idGroup.PUT("", svr.UpdateBook)
		idGroup.DELETE("", svr.DeleteBook)

		idGroup.GET("/items", svr.GetItemsOfBook)
		idGroup.POST("/items", svr.AddItemsToBook)
		idGroup.DELETE("/items", svr.RemoveItemsFromBook)
	}
}
