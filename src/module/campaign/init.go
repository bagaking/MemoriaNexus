package campaign

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

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
	campaignsDetailGroup := group.Group("/campaigns/:id").Use(utils.GinMWParseID())
	{
		campaignsDetailGroup.GET("/monsters", svr.GetMonstersOfCampaign)
		campaignsDetailGroup.GET("/practice", svr.GetMonstersForCampaignPractice)
		campaignsDetailGroup.POST("/submit", svr.SubmitCampaignResult)

		campaignsDetailGroup.GET("/conclusion/today", svr.GetCampaignDungeonConclusionOfToday)
	}
}
