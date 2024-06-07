package profile

import (
	"errors"
	"net/http"

	"github.com/bagaking/memorianexus/internal/utils"

	"github.com/bagaking/goulp/wlog"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"

	"github.com/bagaking/memorianexus/src/model"
	"github.com/bagaking/memorianexus/src/module/dto"
)

// ReqUpdateProfile defines the request format for the UpdateUserProfile endpoint.
type ReqUpdateProfile struct {
	Nickname  string `json:"nickname,omitempty"`
	Email     string `json:"email,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	Bio       string `json:"bio,omitempty"`
}

// GetUserProfile handles a request to retrieve a user's profile information.
// @Summary Get the current user's profile
// @Description Retrieves the profile information for the user who made the request.
// @TagNames profile
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
	profile, err := model.EnsureLoadProfile(svr.db, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.GinHandleError(c, log, http.StatusNotFound, err, "Profile not found")
			return
		}
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
// @TagNames profile
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
	if err := svr.db.Model(updater).Where("id = ?", userID).Updates(updater).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "failed to update profile")
		return
	}

	new(dto.RespProfile).With(new(dto.Profile).FromModel(updater)).Response(c, "profile updated")
}
