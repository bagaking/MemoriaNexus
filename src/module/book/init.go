package book

import (
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
	group.GET("/:id", svr.GetBook)
	group.PUT("/:id", svr.UpdateBook)
	group.DELETE("/:id", svr.DeleteBook)
}
