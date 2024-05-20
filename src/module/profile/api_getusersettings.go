package profile

import (
	"net/http"

	"github.com/bagaking/memorianexus/internal/util"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/gin-gonic/gin"
)

// RespSettingsMemorization defines the structure for the user settings response.
type RespSettingsMemorization struct {
	// Definitions should match with ProfileMemorizationSetting
	ReviewInterval       uint   `json:"review_interval"`
	DifficultyPreference uint8  `json:"difficulty_preference"`
	QuizMode             string `json:"quiz_mode"`
}

// RespSettingsAdvance defines the structure for the advanced settings response.
type RespSettingsAdvance struct {
	Theme              string `json:"theme"`
	Language           string `json:"language"`
	EmailNotifications bool   `json:"email_notifications"`
	PushNotifications  bool   `json:"push_notifications"`
}

// GetUserSettingsMemorization handles a request to get the current user's settings.
// @Summary Get user settings
// @Description Retrieves settings information for the user who made the request.
// @Tags profile
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} RespSettingsMemorization "Successfully retrieved user settings"
// @Failure 400 {object} module.ErrorResponse "Bad Request"
// @Failure 404 {object} module.ErrorResponse "Not Found"
// @Failure 500 {object} module.ErrorResponse "Internal Server Error"
// @Router /profile/settings/memorization [get]
func (svr *Service) GetUserSettingsMemorization(c *gin.Context) {
	userID, exists := util.GetUIDFromGinCtx(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	profile, err := model.EnsureLoadProfile(svr.db, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	settings, err := profile.EnsureLoadProfileSettingsMemorization(svr.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve profile settings"})
		return
	}

	resp := RespSettingsMemorization{
		ReviewInterval:       settings.ReviewInterval,
		DifficultyPreference: settings.DifficultyPreference,
		QuizMode:             settings.QuizMode,
	}

	c.JSON(http.StatusOK, resp)
}

// GetUserSettingsAdvance retrieves advanced settings for the authenticated user.
// @Summary Get user advanced settings
// @Description Retrieves advanced settings information for the current user.
// @Tags profile
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} RespSettingsAdvance "Successfully retrieved user advanced settings"
// @Failure 400 {object} module.ErrorResponse "Bad Request"
// @Failure 404 {object} module.ErrorResponse "Not Found"
// @Failure 500 {object} module.ErrorResponse "Internal Server Error"
// @Router /profile/settings/advance [get]
func (svr *Service) GetUserSettingsAdvance(c *gin.Context) {
	userID, exists := util.GetUIDFromGinCtx(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	profile, err := model.EnsureLoadProfile(svr.db, userID)
	if err != nil || profile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	// Assuming EnsureLoadProfileSettingsAdvance will either load or create if not exists.
	advanceSettings, err := profile.EnsureLoadProfileSettingsAdvance(svr.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve advanced settings"})
		return
	}

	resp := RespSettingsAdvance{
		Theme:              advanceSettings.Theme,
		Language:           advanceSettings.Language,
		EmailNotifications: advanceSettings.EmailNotifications,
		PushNotifications:  advanceSettings.PushNotifications,
	}

	c.JSON(http.StatusOK, resp)
}
