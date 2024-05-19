package book

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Service struct {
	// db *model.db
}

var svr *Service

func Init(db *gorm.DB) (*Service, error) {
	svr = &Service{
		// db: model.NewRepo(db),
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

func (svr *Service) CreateBook(context *gin.Context) {
}

func (svr *Service) GetBooks(context *gin.Context) {
}

func (svr *Service) GetBook(context *gin.Context) {
}

func (svr *Service) UpdateBook(context *gin.Context) {
}

func (svr *Service) DeleteBook(context *gin.Context) {
}
