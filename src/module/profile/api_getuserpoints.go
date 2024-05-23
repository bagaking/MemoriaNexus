package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/model"
)

// RespGetPoints defines the structure for the user profile API response.
type RespGetPoints struct {
	Cash     utils.UInt64 `json:"cash"`
	Gem      utils.UInt64 `json:"gem"`
	VIPScore utils.UInt64 `json:"vip_score"`
}

// GetUserPoints retrieves the points for the authenticated user.
// @Summary Get user points
// @Description Retrieves points information for the current user.
// @Tags profile
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} RespGetPoints "Successfully retrieved user points"
// @Failure 400 {object} module.ErrorResponse "Bad Request"
// @Failure 404 {object} module.ErrorResponse "Not Found"
// @Failure 500 {object} module.ErrorResponse "Internal Server Error"
// @Router /profile/points [get]
func (svr *Service) GetUserPoints(c *gin.Context) {
	userID, exists := utils.GetUIDFromGinCtx(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Assuming that the user points are a part of the Profile model, you would load the profile
	// from the database or service and then extract the points part to respond.
	// This is simulation, replace it with actual logic as needed.
	points, err := model.EnsureLoadProfilePoints(svr.db, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	// Assuming the points information is stored within profile structure.
	resp := RespGetPoints{
		Cash:     points.Cash,
		Gem:      points.Gem,
		VIPScore: points.VipScore,
	}

	c.JSON(http.StatusOK, resp)
}
