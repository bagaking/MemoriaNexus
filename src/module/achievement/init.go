package achievement

import (
	"github.com/bagaking/memorianexus/internal/utils"
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
	group.GET("/", svr.GetAllAchievements)
	idGroup := group.Group("/:id").Use(utils.GinMWParseID())
	{
		idGroup.GET("", svr.GetAchievementDetails)
	}
}

func (svr *Service) GetAllAchievements(context *gin.Context) {
}

func (svr *Service) GetAchievementDetails(context *gin.Context) {
}
