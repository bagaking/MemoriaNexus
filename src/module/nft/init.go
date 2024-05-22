package nft

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

// ApplyMux
//
// NFT 查询
//   - GET /api/v1/nft/nfts：获取用户 NFT
//   - GET /api/v1/nft/nfts/:id：获取 NFT 详情
//
// NFT 操作
//   - POST /api/v1/nft/draw_card：以抽卡的方式创建 nft
//   - POST /api/v1/nft/transfer：赠予
//
// NFT 商店
//   - GET /api/v1/nft/shops：查看所有商店
//   - GET /api/v1/nft/shops/:id：查看某个商店
//
// NFT 交易管理
//   - GET /api/v1/nft/trades：获取市场交易对
//   - POST /api/v1/nft/trades：创建交易对
//   - GET /api/v1/nft/trades/:id：获取交易详情
//   - DELETE /api/v1/nft/trades/:id：取消交易
//   - POST /api/v1/nft/trades/:id/buy：创建购买订单 (会直接生效、因此就是购买)
func (svr *Service) ApplyMux(group gin.IRouter) {
	group.GET("/nfts", svr.GetNFTs)
	group.GET("/nfts/:id", svr.GetNFTDetails)

	group.POST("/draw_card", svr.DrawCard)
	group.POST("/transfer", svr.Transfer)

	group.GET("/shops", svr.GetShops)
	group.GET("/shops/:id", svr.GetShop)

	group.GET("/trades", svr.GetTrades)
	group.POST("/trades", svr.CreateTrade)
	group.GET("/trades/:id", svr.GetTradeDetails)
	group.DELETE("/trades/:id", svr.CancelTrade)
	group.POST("/trades/:id/buy", svr.BuyTrade)
}

func (svr *Service) GetNFTs(c *gin.Context) {
}

func (svr *Service) GetNFTDetails(c *gin.Context) {
}

func (svr *Service) GetShops(c *gin.Context) {
}

func (svr *Service) GetShop(c *gin.Context) {
}

func (svr *Service) DrawCard(c *gin.Context) {
}

func (svr *Service) Transfer(c *gin.Context) {
}

func (svr *Service) GetTrades(c *gin.Context) {
}

func (svr *Service) CreateTrade(c *gin.Context) {
}

func (svr *Service) GetTradeDetails(c *gin.Context) {
}

func (svr *Service) CancelTrade(c *gin.Context) {
}

func (svr *Service) BuyTrade(c *gin.Context) {
}
