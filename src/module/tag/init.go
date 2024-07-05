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
	tagGroup := group.Group("/:tag").Use(utils.GinMWParseTAG())
	{
		tagGroup.GET("/books", svr.GetBooksByTag)
		tagGroup.GET("/items", svr.GetItemsByTag)
	}
}
