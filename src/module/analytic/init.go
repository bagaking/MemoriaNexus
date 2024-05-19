package analytic

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Service struct {
	// db *model.db
}

var svr *Service

func Init(db *gorm.DB) (*Service, error) {
	svr = &Service{
		// db: model.NewRepo(db),
	}
	return svr, nil
}

func (svr *Service) ApplyMux(group gin.IRouter) {
	group.GET("/studyPatterns", svr.GetStudyPatterns)
	group.GET("/timeSpent", svr.GetTimeSpent)
}

func (svr *Service) GetStudyPatterns(context *gin.Context) {
}

func (svr *Service) GetTimeSpent(context *gin.Context) {
}
