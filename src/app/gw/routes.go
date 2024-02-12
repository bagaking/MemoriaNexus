// File: src/app/gw/routes.go

package gw

import (
	"github.com/bagaking/memorianexus/pkg/auth"
	"github.com/bagaking/memorianexus/src/profile/passport"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes - routers all in one
// todo: using rpc
func RegisterRoutes(router gin.IRouter, db *gorm.DB) {

	// todo: 这些值应该从配置中安全获取，现在 MVP 一下
	jwtService := auth.NewJWTService("my_secret_key", "MemoriaNexus")

	svrPassport, _ := passport.Init(db, jwtService)
	svrPassport.ApplyMux(router.Group("/auth"))
}
