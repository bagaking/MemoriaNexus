package trade

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
	group.GET("", svr.GetMarketTrades)
	group.POST("", svr.CreateTrade)
	group.GET("/:id", svr.GetTradeDetails)
	group.DELETE("/:id", svr.CancelTrade)
	group.POST("/:id/purchase", svr.PurchaseTrade)
}

func (svr *Service) GetMarketTrades(context *gin.Context) {
}

func (svr *Service) CreateTrade(context *gin.Context) {
}

func (svr *Service) GetTradeDetails(context *gin.Context) {
}

func (svr *Service) CancelTrade(context *gin.Context) {
}

func (svr *Service) PurchaseTrade(context *gin.Context) {
}
