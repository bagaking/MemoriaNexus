package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/khgame/ranger_iam/pkg/authcli"

	"github.com/bagaking/memorianexus/src/model"
)

// ReqUpdateProfile defines the request format for the UpdateUserProfile endpoint.
type ReqUpdateProfile struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

// UpdateUserProfile handles a request to update the current user's profile information.
// @Summary Update user profile
// @Description Updates the profile information for the user who made the request.
// @Tags profile
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "带有 Bearer 的 Token"
// @Param profile body ReqUpdateProfile true "User profile update info"
// @Success 200 {object} SuccessResponse "Successfully updated user profile"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Not Found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /profile/me [put]
func (svr *Service) UpdateUserProfile(c *gin.Context) {
	// Extract the user ID from the context.
	userID, exists := authcli.GetUIDFromGinCtx(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Bind the incoming JSON to a struct.
	var updateReq ReqUpdateProfile
	if err := c.ShouldBindWith(&updateReq, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	profile := &model.Profile{
		ID:       userID,
		Nickname: updateReq.Nickname,
		Email:    updateReq.Email,
	}

	// Perform the update operation in the repository.
	if err := profile.UpdateProfile(svr.db, profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	// Respond with a generic success message.
	c.JSON(http.StatusOK, SuccessResponse{Message: "Profile updated successfully"})
}
