// File: src/app/gw/routes.go

package gw

import (
	"github.com/gin-gonic/gin"
	"github.com/khgame/ranger_iam/pkg/authcli"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"github.com/bagaking/memorianexus/src/module/achievement"
	"github.com/bagaking/memorianexus/src/module/analytic"
	"github.com/bagaking/memorianexus/src/module/book"
	"github.com/bagaking/memorianexus/src/module/campaign"
	"github.com/bagaking/memorianexus/src/module/dungeon"
	"github.com/bagaking/memorianexus/src/module/item"
	"github.com/bagaking/memorianexus/src/module/nft"
	"github.com/bagaking/memorianexus/src/module/operation"
	"github.com/bagaking/memorianexus/src/module/profile"
	"github.com/bagaking/memorianexus/src/module/system"
	"github.com/bagaking/memorianexus/src/module/tag"
)

//
//func GinMW(iamCli *authcli.Cli) gin.HandlerFunc {
//	AuthN := func(ctx context.Context, tokenStr string) (uint64, error) {
//		uid, err := iamCli.ValidateRemote(ctx, tokenStr)
//		if err == nil {
//			return uid, nil
//		}
//		return 0, err
//	}
//
//	ValidateRemote := func(ctx context.Context, token string) (uint64, error) {
//		// Try to contact the IAM server
//		response, err := cli.httpClient.Get(cli.AuthNSvrURL + "api/v1/session/validate")
//		if err != nil || response.StatusCode != http.StatusOK {
//			// Server is down or returned a non-ok status - use local validation
//			return 0, ErrValidateRemoteStatusFailed
//		}
//
//		// Check if the response includes a degradation signal
//		// Let's assume the header X-Degraded-Mode indicates the mode
//		if response.Header.Get(utils.KEYDegradedMode) == utils.DegradedModeAll {
//			return 0, ErrValidateRemoteDegraded
//		}
//
//		// Server response can be used to get userID
//		defer response.Body.Close()
//		var res struct {
//			UID uint64 `json:"uid"`
//		}
//		if err = json.NewDecoder(response.Body).Decode(&res); err != nil {
//			return 0, err
//		}
//		return res.UID, nil
//	}
//
//	return func(c *gin.Context) {
//		// 从 Header 中获取 tokenStr
//		tokenStr, err := auth.GetTokenStrFromHeader(c)
//		if err != nil {
//			// todo: 从 Cookie 中获取 tokenStr
//			wlog.ByCtx(c).WithError(err).Error("authn failed, tokenStr is empty")
//			c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
//			return
//		}
//
//		uid, err := AuthN(c, tokenStr)
//		if err != nil {
//			wlog.ByCtx(c).WithError(err).Errorf("authn failed, uid is empty, tokenStr= %s", tokenStr)
//			c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
//			return
//		}
//		c.Set(authcli.UserCtxKey, uid)
//		c.Next()
//	}
//}

func RegisterCallbacks(router gin.IRouter) {
	group := router.Group("/lark_openapi")
	group.GET("/event", LarkEventHandler)
}

// RegisterRoutes - routers all in one
// todo: using rpc
func RegisterRoutes(router gin.IRouter, db *gorm.DB, iamCli *authcli.Cli) {
	// 用户账户服务路由组
	svrProfile := profile.NewService(db)
	g := router.Group("/profile")
	g.Use(iamCli.GinMW())
	svrProfile.ApplyMux(g)

	// 学习材料管理路由组
	svrItems := item.NewService(db)
	g = router.Group("/items")
	g.Use(iamCli.GinMW())
	svrItems.ApplyMux(g)

	// 册子管理路由组
	svrBooks, _ := book.Init(db)
	g = router.Group("/books")
	g.Use(iamCli.GinMW())
	svrBooks.ApplyMux(g)

	// 标签管理路由组
	svrTags, _ := tag.Init(db)
	g = router.Group("/tags")
	g.Use(iamCli.GinMW())
	svrTags.ApplyMux(g)

	// 系统操作路由组
	svrSystem, _ := system.Init(db)
	g = router.Group("/system")
	g.Use(iamCli.GinMW())
	svrSystem.ApplyMux(g)

	// 复习计划管理路由组
	svrDungeon, _ := dungeon.Init(db)
	g = router.Group("/dungeon")
	g.Use(iamCli.GinMW())
	svrDungeon.ApplyMux(g)

	svrCampaignDungeon, _ := campaign.Init(db)
	g = router.Group("/dungeon")
	g.Use(iamCli.GinMW())
	svrCampaignDungeon.ApplyMux(g)

	// 数据分析路由组
	svrAnalytics, _ := analytic.Init(db)
	g = router.Group("/analytic")
	g.Use(iamCli.GinMW())
	svrAnalytics.ApplyMux(g)

	// NFT管理路由组
	svrNfts, _ := nft.Init(db)
	g = router.Group("/nft")
	g.Use(iamCli.GinMW())
	svrNfts.ApplyMux(g)

	// 成就系统路由组
	svrAchievements, _ := achievement.Init(db)
	g = router.Group("/achievements")
	g.Use(iamCli.GinMW())
	svrAchievements.ApplyMux(g)

	// 运营管理路由组
	svrOperation, _ := operation.Init(db)
	g = router.Group("/operation")
	g.Use(iamCli.GinMW())
	svrOperation.ApplyMux(g)

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
