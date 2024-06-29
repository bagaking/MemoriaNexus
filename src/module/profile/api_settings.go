package profile

import (
	"net/http"

	"github.com/bagaking/memorianexus/internal/utils"

	"github.com/bagaking/goulp/wlog"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/bagaking/memorianexus/src/model"
	"github.com/bagaking/memorianexus/src/module/dto"
)

// GetUserSettingsMemorization handles a request to get the current user's settings.
// @Summary Get user settings
// @Description Retrieves settings information for the user who made the request.
// @Tags profile
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} dto.RespSettingsMemorization "Successfully retrieved user settings"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 404 {object} utils.ErrorResponse "Not Found"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /profile/settings/memorization [get]
func (svr *Service) GetUserSettingsMemorization(c *gin.Context) {
	userID := utils.GinMustGetUserID(c)
	log := wlog.ByCtx(c, "GetUserSettingsMemorization").WithField("user_id", userID)

	profile, err := model.EnsureLoadProfile(svr.db, userID)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusNotFound, err, "Profile not found")
		return
	}

	settings, err := profile.EnsureLoadProfileSettingsMemorization(svr.db)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to retrieve profile settings")
		return
	}

	new(dto.RespSettingsMemorization).With(new(dto.SettingsMemorization).FromModel(settings)).Response(c)
}

// GetUserSettingsAdvance retrieves advanced settings for the authenticated user.
// @Summary Get user advanced settings
// @Description Retrieves advanced settings information for the current user.
// @Tags profile
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} dto.RespSettingsAdvance "Successfully retrieved user advanced settings"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 404 {object} utils.ErrorResponse "Not Found"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /profile/settings/advance [get]
func (svr *Service) GetUserSettingsAdvance(c *gin.Context) {
	userID := utils.GinMustGetUserID(c)
	log := wlog.ByCtx(c, "GetUserSettingsAdvance").WithField("user_id", userID)

	profile, err := model.EnsureLoadProfile(svr.db, userID)
	if err != nil || profile == nil {
		utils.GinHandleError(c, log, http.StatusNotFound, err, "Profile not found")
		return
	}

	// Assuming EnsureLoadProfileSettingsAdvance will either load or create if not exists.
	advanceSettings, err := profile.EnsureLoadProfileSettingsAdvance(svr.db)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to retrieve advanced settings")
		return
	}

	new(dto.RespSettingsAdvance).With(new(dto.SettingsAdvance).FromModel(advanceSettings)).Response(c)
}

// UpdateUserSettingsMemorization handles a request to update the current user's settings.
// @Summary Update user settings
// @Description Updates the settings for the user who made the request.
// @Tags profile
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param settings body ReqUpdateUserSettingsMemorization true "User settings update info"
// @Success 200 {object} dto.RespSettingsMemorization "Successfully updated user settings"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 404 {object} utils.ErrorResponse "Not Found"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /profile/settings/memorization [put]
func (svr *Service) UpdateUserSettingsMemorization(c *gin.Context) {
	userID := utils.GinMustGetUserID(c)
	log := wlog.ByCtx(c, "UpdateUserSettingsMemorization").WithField("user_id", userID)

	var updateReq ReqUpdateUserSettingsMemorization
	if err := c.ShouldBindWith(&updateReq, binding.JSON); err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	profile, err := model.EnsureLoadProfile(svr.db, userID)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusNotFound, err, "Profile not found")
		return
	}

	settingsToUpdate := &model.ProfileMemorizationSetting{
		ID: userID,
	}

	// Update the fields that were provided in the request.
	if updateReq.ReviewIntervalSetting != nil {
		settingsToUpdate.ReviewIntervalSetting = *updateReq.ReviewIntervalSetting
	}
	if updateReq.DifficultyPreference != nil {
		settingsToUpdate.DifficultyPreference = *updateReq.DifficultyPreference
	}
	if updateReq.QuizMode != nil {
		settingsToUpdate.QuizMode = *updateReq.QuizMode
	}

	if err = profile.UpdateSettingsMemorization(svr.db, settingsToUpdate); err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to update settings")
		return
	}

	new(dto.RespSettingsMemorization).With(new(dto.SettingsMemorization).FromModel(settingsToUpdate)).
		Response(c, "memorization settings updated")
}

// UpdateUserSettingsAdvance updates the advanced settings for the current user.
// @Summary Update user advanced settings
// @Description Updates advanced settings for the authenticated user.
// @Tags profile
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param settings body ReqUpdateUserSettingsAdvance true "User advanced settings update info"
// @Success 200 {object} dto.RespSettingsAdvance "Successfully updated user advanced settings"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 404 {object} utils.ErrorResponse "Not Found"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /profile/settings/advance [put]
func (svr *Service) UpdateUserSettingsAdvance(c *gin.Context) {
	userID := utils.GinMustGetUserID(c)
	log := wlog.ByCtx(c, "UpdateUserSettingsAdvance").WithField("user_id", userID)

	var updateReq ReqUpdateUserSettingsAdvance
	if err := c.ShouldBindWith(&updateReq, binding.JSON); err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	profile, err := model.EnsureLoadProfile(svr.db, userID)
	if err != nil || profile == nil {
		utils.GinHandleError(c, log, http.StatusNotFound, err, "Profile not found")
		return
	}

	advanceSettings := &model.ProfileAdvanceSetting{
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

	if err = profile.UpdateSettingsAdvance(svr.db, advanceSettings); err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to update advanced settings")
		return
	}

	new(dto.RespSettingsAdvance).With(
		new(dto.SettingsAdvance).FromModel(advanceSettings),
	).Response(c, "advanced settings updated")
}
