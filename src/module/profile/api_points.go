package profile

import (
	"net/http"

	"github.com/bagaking/goulp/wlog"
	"github.com/gin-gonic/gin"

	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/bagaking/memorianexus/src/module/dto"
)

// GetUserPoints retrieves the points for the authenticated user.
// @Summary Get user points
// @Description Retrieves points information for the current user.
// @Tags profile
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} dto.RespPoints "Successfully retrieved user points"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 404 {object} utils.ErrorResponse "Not Found"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /profile/points [get]
func (svr *Service) GetUserPoints(c *gin.Context) {
	userID := utils.GinMustGetUserID(c)
	log := wlog.ByCtx(c, "GetUserPoints").WithField("user_id", userID)

	// Assuming that the user points are a part of the Profile model, you would load the profile
	// from the database or service and then extract the points part to respond.
	// This is simulation, replace it with actual logic as needed.
	points, err := model.EnsureLoadProfilePoints(svr.db, userID)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusNotFound, err, "Profile not found")
		return
	}

	// Assuming the points information is stored within profile structure.
	new(dto.RespPoints).With(new(dto.Points).FromModel(points)).Response(c)
}
