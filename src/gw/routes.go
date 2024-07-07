// File: src/app/gw/routes.go

package gw

import (
	"github.com/bagaking/memorianexus/src/module/campaign"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"github.com/khgame/ranger_iam/pkg/authcli"

	"github.com/bagaking/memorianexus/src/module/achievement"
	"github.com/bagaking/memorianexus/src/module/analytic"
	"github.com/bagaking/memorianexus/src/module/book"
	"github.com/bagaking/memorianexus/src/module/dungeon"
	"github.com/bagaking/memorianexus/src/module/item"
	"github.com/bagaking/memorianexus/src/module/nft"
	"github.com/bagaking/memorianexus/src/module/operation"
	"github.com/bagaking/memorianexus/src/module/profile"
	"github.com/bagaking/memorianexus/src/module/system"
	"github.com/bagaking/memorianexus/src/module/tag"
)

// RegisterRoutes - routers all in one
// todo: using rpc
func RegisterRoutes(router gin.IRouter, db *gorm.DB) {
	// todo: 这些值应该从配置中安全获取，现在 MVP 一下
	iamCli := authcli.New("my_secret_key", "http://0.0.0.0/")

	router.Use(iamCli.GinMW())

	// 用户账户服务路由组
	svrProfile := profile.NewService(db)
	svrProfile.ApplyMux(router.Group("/profile"))

	// 学习材料管理路由组
	svrItems := item.NewService(db)
	svrItems.ApplyMux(router.Group("/items"))

	// 册子管理路由组
	svrBooks, _ := book.Init(db)
	svrBooks.ApplyMux(router.Group("/books"))

	// 标签管理路由组
	svrTags, _ := tag.Init(db)
	svrTags.ApplyMux(router.Group("/tags"))

	// 系统操作路由组
	svrSystem, _ := system.Init(db)
	svrSystem.ApplyMux(router.Group("/system"))

	// 复习计划管理路由组
	svrDungeon, _ := dungeon.Init(db)
	svrDungeon.ApplyMux(router.Group("/dungeon"))

	svrCampaignDungeon, _ := campaign.Init(db)
	svrCampaignDungeon.ApplyMux(router.Group("/dungeon"))

	// 数据分析路由组
	svrAnalytics, _ := analytic.Init(db)
	svrAnalytics.ApplyMux(router.Group("/analytic"))

	// NFT管理路由组
	svrNfts, _ := nft.Init(db)
	svrNfts.ApplyMux(router.Group("/nft"))

	// 成就系统路由组
	svrAchievements, _ := achievement.Init(db)
	svrAchievements.ApplyMux(router.Group("/achievements"))

	// 运营管理路由组
	svrOperation, _ := operation.Init(db)
	svrOperation.ApplyMux(router.Group("/operation"))

	// 社区互动路由组
	//svrCommunity, _ := community.Init(db)
	//svrCommunity.ApplyMux(router.Group("/community"))
	//{
	//	group.POST("/post", svr.CreatePost)
	//	group.GET("/post", svr.GetPosts)
	//	group.PUT("/post/:id", svr.UpdatePost)
	//	group.DELETE("/post/:id", svr.DeletePost)
	//	group.POST("/post/:id/comments", svr.CommentPost)
	//	group.GET("/post/:id/comments", svr.GetPostComments)
	//}

	// 注册 Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
