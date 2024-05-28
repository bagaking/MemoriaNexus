package dungeon

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

	// 自动迁移数据库结构
	//err := db.AutoMigrate(&model.Dungeon{}, &model.DungeonBook{}, &model.DungeonItem{})
	//if err != nil {
	//	return nil, err
	//}

	return svr, nil
}

func (svr *Service) ApplyMux(group gin.IRouter) {
	group.POST("/dungeons", svr.CreateDungeon)
	group.GET("/dungeons", svr.GetDungeons)

	group.GET("/dungeons/:id", svr.GetDungeon)
	group.DELETE("/dungeons/:id", svr.DeleteDungeon)

	group.PUT("/dungeons/:id", svr.UpdateDungeon)
	group.POST("/dungeons/:id/books", svr.AppendBooksToDungeon)
	group.POST("/dungeons/:id/items", svr.AppendItemsToDungeon)
	group.POST("/dungeons/:id/tags", svr.AppendTagsToDungeon)

	group.GET("/dungeons/:id/books", svr.GetDungeonBooksDetail)
	group.GET("/dungeons/:id/items", svr.GetDungeonItemsDetail)
	group.GET("/dungeons/:id/tags", svr.GetDungeonTagsDetail)

	group.DELETE("/dungeons/:id/books", svr.SubtractDungeonBooks)
	group.DELETE("/dungeons/:id/items", svr.SubtractDungeonItems)
	group.DELETE("/dungeons/:id/tags", svr.SubtractDungeonTags)

	// 新增的复习相关API接口
	group.GET("/campaigns/:id/monsters", svr.GetMonstersOfCampaignDungeon)
	group.GET("/campaigns/:id/next_monsters", svr.GetNextMonstersOfCampaignDungeon)
	group.GET("/campaigns/:id/today_conclusion", svr.GetCampaignDungeonTodayConclusion)
	group.POST("/campaigns/:id/report_result", svr.ReportCampaignResult)

	group.GET("/endless/:id/monsters", svr.GetMonstersOfEndlessDungeon)
	// group.GET("/endless/:id/next_monsters", svr.GetNextMonstersOfEndlessDungeon)
	// group.GET("/endless/:id/today_conclusion", svr.GetEndlessDungeonTodayConclusion)
	// group.POST("/endless/:id/report_result", svr.ReportEndlessResult)

	group.GET("/instances", svr.GetDungeonInstances)
	group.GET("/instances/:id", svr.GetDungeonInstance)
}

func (svr *Service) GetDungeonInstances(context *gin.Context) {
	// 实现代码
}

func (svr *Service) GetDungeonInstance(context *gin.Context) {
	// 实现代码
}
