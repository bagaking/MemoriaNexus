package dungeon

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
	group.GET("/schedules", svr.GetDungeonSchedules)
	group.GET("/schedules/:id", svr.GetDungeonSchedule)
	group.POST("/schedules", svr.CreateDungeonSchedule)
	group.PUT("/schedules/:id", svr.UpdateDungeonSchedule)
	group.DELETE("/schedules/:id", svr.DeleteDungeonSchedule)
	group.GET("/instances", svr.GetDungeonInstances)
	group.GET("/instances/:id", svr.GetDungeonInstance)
}

func (svr *Service) GetDungeonSchedules(context *gin.Context) {
}

func (svr *Service) GetDungeonSchedule(context *gin.Context) {
}

func (svr *Service) CreateDungeonSchedule(context *gin.Context) {
}

func (svr *Service) UpdateDungeonSchedule(context *gin.Context) {
}

func (svr *Service) DeleteDungeonSchedule(context *gin.Context) {
}

func (svr *Service) GetDungeonInstances(context *gin.Context) {
}

func (svr *Service) GetDungeonInstance(context *gin.Context) {
}
