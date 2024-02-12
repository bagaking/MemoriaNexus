package passport

import (
	"github.com/bagaking/memorianexus/pkg/auth"
	"github.com/bagaking/memorianexus/src/iam/passport/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Service struct {
	Repo *model.Repo
	JWT  *auth.JWTService
}

var svr *Service

func Init(db *gorm.DB, jwtService *auth.JWTService) (*Service, error) {
	svr = &Service{
		Repo: model.NewRepo(db),
		JWT:  jwtService,
	}
	return svr, nil
}

func (svr *Service) ApplyMux(group gin.IRouter) {
	group.POST("/register", svr.HandleRegister)
	group.POST("/login", svr.HandleLogin)
}
