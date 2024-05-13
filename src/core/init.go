package core

import (
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

	"github.com/bagaking/memorianexus/src/core/model"
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
	// todo
}
