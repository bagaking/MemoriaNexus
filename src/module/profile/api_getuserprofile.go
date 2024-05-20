package profile

import (
	"errors"
	"net/http"

	"github.com/bagaking/memorianexus/internal/util"

	"github.com/bagaking/goulp/wlog"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RespGetProfile defines the structure for the user profile API response.
type RespGetProfile struct {
	ID        util.UInt64 `json:"id"`
	Nickname  string      `json:"nickname"`
	Email     string      `json:"email"`
	AvatarURL string      `json:"avatar_url"`
	// Include other fields as appropriate.
}

// mapProfileToResponse is a helper function to map the Profile model to the API response struct.
// (This implementation assumes you have a separate struct for your API responses.
//
//	You would need to implement this mapping according to your actual API response structure.)
func mapProfileToResponse(profile *model.Profile) *RespGetProfile {
	return &RespGetProfile{
		ID:        profile.ID,
		Nickname:  profile.Nickname,
		Email:     profile.Email,
		AvatarURL: profile.AvatarURL,
		// Additional fields should be mapped as required.
	}
}

// GetUserProfile handles a request to retrieve a user's profile information.
// @Summary Get the current user's profile
// @Description Retrieves the profile information for the user who made the request.
// @Tags profile
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} RespGetProfile "Successfully retrieved user profile"
// @Failure 400 {object} module.ErrorResponse "Bad Request"
// @Failure 404 {object} module.ErrorResponse "Not Found"
// @Failure 500 {object} module.ErrorResponse "Internal Server Error"
// @Router /profile/me [get]
func (svr *Service) GetUserProfile(c *gin.Context) {
	log := wlog.ByCtx(c)
	// Extract the user ID from the context.
	userID, exists := util.GetUIDFromGinCtx(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Use the ID to load the profile from the database.
	profile, err := model.EnsureLoadProfile(svr.db, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Profile not found", "error": err.Error()})
			log.WithError(err).Warnf("load profile by uid %v not found", userID)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving profile", "error": err.Error()})
		log.WithError(err).Errorf("load profile by uid %v failed", userID)
		return
	}

	// Map the model to the output format.
	profileResponse := mapProfileToResponse(profile)

	// Respond with the user profile data.
	c.JSON(http.StatusOK, profileResponse)
}
