package system

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
	group.GET("/notifications", svr.GetAllNotifications)
	group.POST("/notifications/markAsRead", svr.MarkNotificationsAsRead)
	group.GET("/announcements", svr.GetAllAnnouncements)
	group.POST("/announcements/markAsRead", svr.MarkAnnouncementsAsRead)
	group.GET("/configs", svr.GetGlobalConfigs)
}

func (svr *Service) GetAllNotifications(context *gin.Context) {
}

func (svr *Service) MarkNotificationsAsRead(context *gin.Context) {
}

func (svr *Service) GetAllAnnouncements(context *gin.Context) {
}

func (svr *Service) MarkAnnouncementsAsRead(context *gin.Context) {
}

func (svr *Service) GetGlobalConfigs(context *gin.Context) {
}
