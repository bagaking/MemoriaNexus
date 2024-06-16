package tag

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
	group.GET("", svr.GetTags)
	idGroup := group.Group("/:id").Use(utils.GinMWParseID())
	{
		idGroup.GET("/items", svr.GetItemsByTag)
		idGroup.POST("/books", svr.GetBooksByTag)
	}

	group.GET("/name/:name", svr.GetTagByName)
	group.GET("/name/:name/books", svr.GetBooksByTagName)
	group.GET("/name/:name/items", svr.GetItemsByTagName)
}
