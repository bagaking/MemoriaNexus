package operation

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
	group.GET("/task", svr.GetCurrentTasks)
	group.GET("/task/:id", svr.GetTaskDetails)
	group.GET("/activity", svr.GetActivities)
	group.GET("/activity/:id", svr.GetActivityDetails)
}

func (svr *Service) GetCurrentTasks(context *gin.Context) {
}

func (svr *Service) GetTaskDetails(context *gin.Context) {
}

func (svr *Service) GetActivities(context *gin.Context) {
}

func (svr *Service) GetActivityDetails(context *gin.Context) {
}
