package dungeon

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

	// 自动迁移数据库结构
	//err := db.AutoMigrate(&model.Dungeon{}, &model.DungeonBook{}, &model.DungeonItem{})
	//if err != nil {
	//	return nil, err
	//}

	return svr, nil
}

func (svr *Service) ApplyMux(group gin.IRouter) {
	dungeonsGeneralGroup := group.Group("/dungeons")
	{
		dungeonsGeneralGroup.POST("", svr.CreateDungeon)
		dungeonsGeneralGroup.GET("", svr.GetDungeons)

		dungeonsDetailGroup := dungeonsGeneralGroup.Group("/:id").Use(utils.GinMWParseID())
		{
			dungeonsDetailGroup.GET("", svr.GetDungeon)
			dungeonsDetailGroup.DELETE("", svr.DeleteDungeon)
			dungeonsDetailGroup.PUT("", svr.UpdateDungeon)

			dungeonsDetailGroup.POST("/books", svr.AppendBooksToDungeon)
			dungeonsDetailGroup.POST("/items", svr.AppendItemsToDungeon)
			dungeonsDetailGroup.POST("/tags", svr.AppendTagsToDungeon)

			dungeonsDetailGroup.GET("/books", svr.GetDungeonBooksDetail)
			dungeonsDetailGroup.GET("/items", svr.GetDungeonItemsDetail)
			dungeonsDetailGroup.GET("/tags", svr.GetDungeonTagsDetail)

			dungeonsDetailGroup.DELETE("/books", svr.SubtractDungeonBooks)
			dungeonsDetailGroup.DELETE("/items", svr.SubtractDungeonItems)
			dungeonsDetailGroup.DELETE("/tags", svr.SubtractDungeonTags)
		}
	}

	campaignsDetailGroup := group.Group("/campaigns/:id").Use(utils.GinMWParseID())
	{
		campaignsDetailGroup.GET("/monsters", svr.GetMonstersOfCampaignDungeon)
		campaignsDetailGroup.GET("/practice", svr.GetMonstersForCampaignPractice)
		campaignsDetailGroup.POST("/submit", svr.SubmitCampaignResult)

		campaignsDetailGroup.GET("/conclusion/today", svr.GetCampaignDungeonConclusionOfToday)
	}

	endlessDetailGroup := group.Group("/endless/:id").Use(utils.GinMWParseID())
	{
		endlessDetailGroup.GET("/monsters", svr.GetMonstersOfEndlessDungeon)
		// endlessDetailGroup.GET("/next_monsters", svr.GetNextMonstersOfEndlessDungeon)
		// endlessDetailGroup.GET("/today_conclusion", svr.GetEndlessDungeonTodayConclusion)
		// endlessDetailGroup.POST("/report_result", svr.ReportEndlessResult)
	}

	group.GET("/instances/:id", svr.GetDungeonInstance)
}

func (svr *Service) GetDungeonInstance(context *gin.Context) {
	// 实现代码
}
