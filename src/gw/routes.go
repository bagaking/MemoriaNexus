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

// RegisterRoutes - routers all in one
// todo: using rpc
func RegisterRoutes(router gin.IRouter, db *gorm.DB, iamCli *authcli.Cli) {
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
