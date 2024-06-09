package book

import (
	"net/http"

	"github.com/bagaking/goulp/wlog"
	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/bagaking/memorianexus/src/module/dto"
	"github.com/gin-gonic/gin"
)

// ReqGetBookItemsQuery encapsulates the request parameters for fetching items.
type ReqGetBookItemsQuery struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

// GetBookItems handles retrieving a list of books with pagination.
// @Summary Get item list of books with pagination
// @Description Get a paginated list of items for the book.
// @Tags book
// @Accept json
// @Produce json
// @Param page query int false "Page number for pagination" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} dto.RespItemList "items of the book found"
// @Router /books/{id}/items [get]
func (svr *Service) GetBookItems(c *gin.Context) {
	userID := utils.GinMustGetUserID(c)
	bookID := utils.GinMustGetID(c)
	log := wlog.ByCtx(c, "GetBookItems").WithField("user_id", userID)

	var req ReqGetBookItemsQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, err, "Invalid query parameters")
		return
	}
	// Set default values for pagination.
	if req.Page < 1 {
		req.Page = 1
	}

	if req.Limit < 1 {
		req.Limit = 10
	}

	pager := new(dto.RespItemList).SetPageAndLimit(req.Page, req.Limit)
	items, err := model.GetItemsOfBook(svr.db, bookID, pager.Offset, req.Limit)
	if err != nil {
		log.WithError(err).Errorf("Failed to fetch books for user %v", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching book"})
	}

	for _, item := range items {
		pager.Append(new(dto.Item).FromModel(item))
	}
	pager.Response(c, "items of the book found")
}
