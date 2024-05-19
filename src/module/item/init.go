package item

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
	group.POST("", svr.CreateItem)
	group.GET("", svr.GetItems)
	group.GET("/:id", svr.GetItem)
	group.PUT("/:id", svr.UpdateItem)
	group.DELETE("/:id", svr.DeleteItem)
}

func (svr *Service) CreateItem(context *gin.Context) {
}

func (svr *Service) GetItems(context *gin.Context) {
}

func (svr *Service) GetItem(context *gin.Context) {
}

func (svr *Service) UpdateItem(context *gin.Context) {
}

func (svr *Service) DeleteItem(context *gin.Context) {
}
