package profile

import (
	"net/http"

	"github.com/bagaking/memorianexus/internal/utils"

	"github.com/bagaking/goulp/wlog"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/bagaking/memorianexus/src/module/dto"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// GetUserProfile handles a request to retrieve a user's profile information.
// @Summary Get the current user's profile
// @Description Retrieves the profile information for the user who made the request.
// @Tags profile
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.RespProfile "Successfully retrieved user profile"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 404 {object} utils.ErrorResponse "Not Found"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /profile/me [get]
func (svr *Service) GetUserProfile(c *gin.Context) {
	userID := utils.GinMustGetUserID(c)
	log := wlog.ByCtx(c, "GetUserProfile").WithField("user_id", userID)

	// Use the ID to load the profile from the database.
	profile, err := model.EnsureProfile(c, svr.db, userID)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Error retrieving profile")
		return
	}

	// Respond with the user profile data.
	new(dto.RespProfile).With(
		new(dto.Profile).FromModel(profile),
	).Response(c, "profile found")
}

// UpdateUserProfile handles a request to update the current user's profile information.
// @Summary Update user profile
// @Description Updates the profile information for the user who made the request.
// @Tags profile
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param profile body ReqUpdateProfile true "User profile update info"
// @Success 200 {object} dto.SuccessResponse "Successfully updated user profile"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 404 {object} utils.ErrorResponse "Not Found"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /profile/me [put]
func (svr *Service) UpdateUserProfile(c *gin.Context) {
	userID := utils.GinMustGetUserID(c)
	log := wlog.ByCtx(c, "UpdateUserProfile").WithField("user_id", userID)

	var req ReqUpdateProfile
	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, err, "invalid request body")
		return
	}

	updater := &model.Profile{
		ID:        userID,
		Nickname:  req.Nickname,
		Email:     req.Email,
		AvatarURL: req.AvatarURL,
		Bio:       req.Bio,
	}

	// Perform the update operation in the repository.
	// todo: its not work for now
	if err := updater.UpdateProfile(svr.db); err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "failed to update profile")
		return
	}
	log.Infof("profile updated, updater= %v", updater)

	new(dto.RespProfile).With(new(dto.Profile).FromModel(updater)).Response(c, "profile updated")
}
