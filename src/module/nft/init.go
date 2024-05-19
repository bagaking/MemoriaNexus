package nft

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
	group.GET("", svr.GetUserNFTs)
	group.GET("/:id", svr.GetNFTDetails)
	group.GET("/shop", svr.ViewShop)
	group.POST("/draw", svr.DrawCard)
}

func (svr *Service) GetUserNFTs(context *gin.Context) {
}

func (svr *Service) GetNFTDetails(context *gin.Context) {
}

func (svr *Service) ViewShop(context *gin.Context) {
}

func (svr *Service) DrawCard(context *gin.Context) {
}
