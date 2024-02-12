package profile

import (
	"github.com/bagaking/memorianexus/src/profile/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Service struct {
	Repo *model.Repo
}

var svr *Service

func Init(db *gorm.DB) (*Service, error) {
	svr = &Service{
		Repo: model.NewRepo(db),
	}
	return svr, nil
}

func (svr *Service) ApplyMux(group gin.IRouter) {
	group.PUT("/settings", svr.SetProfileSettings)
}
