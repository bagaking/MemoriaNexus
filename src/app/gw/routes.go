// File: src/app/gw/routes.go

package gw

import (
	"github.com/bagaking/memorianexus/pkg/auth"
	"github.com/bagaking/memorianexus/src/iam/passport"
	"github.com/bagaking/memorianexus/src/profile"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes - routers all in one
// todo: using rpc
func RegisterRoutes(router gin.IRouter, db *gorm.DB) {

	// todo: 这些值应该从配置中安全获取，现在 MVP 一下
	jwtService := auth.NewJWTService("my_secret_key", "MemoriaNexus")
	nwAuth := jwtService.AuthMiddleware()

	authGroup := router.Group("/auth")
	{
		svrPassport, _ := passport.Init(db, jwtService)
		svrPassport.ApplyMux(authGroup)
	}

	profileGroup := router.Group("/profile")
	profileGroup.Use(nwAuth)
	{
		svrProfile, _ := profile.Init(db)
		svrProfile.ApplyMux(profileGroup)
	}

	coreGroup := router.Group("/core")
	coreGroup.Use(nwAuth)
	{

	}
}
