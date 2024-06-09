package book

import (
	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	group.GET("", svr.GetBooks)
	idGroup := group.Group("/:id").Use(utils.GinMWParseID())
	{
		idGroup.GET("", svr.GetBook)
		idGroup.PUT("", svr.UpdateBook)
		idGroup.DELETE("", svr.DeleteBook)
		idGroup.GET("/items", svr.GetBookItems)
	}
}
