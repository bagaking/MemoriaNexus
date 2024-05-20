package profile

import (
	"net/http"

	"github.com/bagaking/memorianexus/internal/util"

	"github.com/bagaking/memorianexus/src/module"

	"github.com/bagaking/memorianexus/src/model"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// ReqUpdateUserSettingsMemorization defines the request format for updating user settings.
type ReqUpdateUserSettingsMemorization struct {
	ReviewInterval       *uint   `json:"review_interval"`
	DifficultyPreference *uint8  `json:"difficulty_preference"`
	QuizMode             *string `json:"quiz_mode"`
}

// ReqUpdateUserSettingsAdvance defines the request to update advanced settings.
type ReqUpdateUserSettingsAdvance struct {
	Theme              *string `json:"theme"`
	Language           *string `json:"language"`
	EmailNotifications *bool   `json:"email_notifications"`
	PushNotifications  *bool   `json:"push_notifications"`
}

// UpdateUserSettingsMemorization handles a request to update the current user's settings.
// @Summary Update user settings
// @Description Updates the settings for the user who made the request.
// @Tags profile
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param settings body ReqUpdateUserSettingsMemorization true "User settings update info"
// @Success 200 {object} module.SuccessResponse "Successfully updated user settings"
// @Failure 400 {object} module.ErrorResponse "Bad Request"
// @Failure 404 {object} module.ErrorResponse "Not Found"
// @Failure 500 {object} module.ErrorResponse "Internal Server Error"
// @Router /profile/settings/memorization [put]
func (svr *Service) UpdateUserSettingsMemorization(c *gin.Context) {
	userID, exists := util.GetUIDFromGinCtx(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var updateReq ReqUpdateUserSettingsMemorization
	if err := c.ShouldBindWith(&updateReq, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	profile, err := model.EnsureLoadProfile(svr.db, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	settingsToUpdate := model.ProfileMemorizationSetting{
		ID: userID,
	}

	// Update the fields that were provided in the request.
	if updateReq.ReviewInterval != nil {
		settingsToUpdate.ReviewInterval = *updateReq.ReviewInterval
	}
	if updateReq.DifficultyPreference != nil {
		settingsToUpdate.DifficultyPreference = *updateReq.DifficultyPreference
	}
	if updateReq.QuizMode != nil {
		settingsToUpdate.QuizMode = *updateReq.QuizMode
	}

	err = profile.SaveProfileSettingsMemorization(svr.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
		return
	}

	c.JSON(http.StatusOK, module.SuccessResponse{Message: "Settings updated successfully"})
}

// UpdateUserSettingsAdvance updates the advanced settings for the current user.
// @Summary Update user advanced settings
// @Description Updates advanced settings for the authenticated user.
// @Tags profile
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param settings body ReqUpdateUserSettingsAdvance true "User advanced settings update info"
// @Success 200 {object} module.SuccessResponse "Successfully updated user advanced settings"
// @Failure 400 {object} module.ErrorResponse "Bad Request"
// @Failure 404 {object} module.ErrorResponse "Not Found"
// @Failure 500 {object} module.ErrorResponse "Internal Server Error"
// @Router /profile/settings/advance [put]
func (svr *Service) UpdateUserSettingsAdvance(c *gin.Context) {
	userID, exists := util.GetUIDFromGinCtx(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var updateReq ReqUpdateUserSettingsAdvance
	if err := c.ShouldBindWith(&updateReq, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	profile, err := model.EnsureLoadProfile(svr.db, userID)
	if err != nil || profile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	advanceSettings := model.ProfileAdvanceSetting{
		ID: userID,
	}

	// Update the fields that were provided in the request.
	if updateReq.Theme != nil {
		advanceSettings.Theme = *updateReq.Theme
	}
	if updateReq.Language != nil {
		advanceSettings.Language = *updateReq.Language
	}
	if updateReq.EmailNotifications != nil {
		advanceSettings.EmailNotifications = *updateReq.EmailNotifications
	}
	if updateReq.PushNotifications != nil {
		advanceSettings.PushNotifications = *updateReq.PushNotifications
	}

	err = profile.SaveProfileSettingsAdvance(svr.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update advanced settings"})
		return
	}

	c.JSON(http.StatusOK, module.SuccessResponse{Message: "Advanced settings updated successfully"})
}
