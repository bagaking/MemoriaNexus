// File: src/app/gw/routes.go

package gw

import (
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

	"github.com/bagaking/memorianexus/src/profile/passport"
)

// RegisterRoutes - routers all in one
// todo: using rpc
func RegisterRoutes(router gin.IRouter, db *gorm.DB) {

	svrPassport := passport.NewRegisterService(db)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", svrPassport.HandleRegister)
	}
}
