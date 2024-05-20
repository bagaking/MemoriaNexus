package item

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bagaking/memorianexus/internal/util"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/bagaking/memorianexus/src/module"
)

// GetItemsRequest encapsulates the request parameters for fetching items.
type GetItemsRequest struct {
	UserID util.UInt64 `form:"user_id"`
	BookID util.UInt64 `form:"book_id"`
	Type   string      `form:"type"`
	Page   int         `form:"page"`
	Limit  int         `form:"limit"`
}

// GetItems handles retrieving a list of items with optional filters and pagination.
// @Summary Get a list of items with optional filters
// @Description Get a list of items for the user with optional filters for book and type and support for pagination.
// @Tags item
// @Accept json
// @Produce json
// @Param user_id query uint64 false "User ID"
// @Param book_id query uint64 false "Book ID"
// @Param type query string false "Type of item"
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} []model.Item "Successfully retrieved items"
// @Failure 400 {object} module.ErrorResponse "Bad Request"
// @Router /items [get]
func (svr *Service) GetItems(c *gin.Context) {
	userID, exists := util.GetUIDFromGinCtx(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req GetItemsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, module.ErrorResponse{Message: "Invalid query parameters"})
		return
	}

	// Set default values for pagination.
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}

	query := svr.db.Model(&model.Item{})
	if req.UserID <= 0 { // 如果不指定用户，搜索的就是自己的 todo：要用这个接口支持搜索其他人的吗？
		req.UserID = userID
	}
	query = query.Where("user_id = ?", req.UserID)
	if req.BookID > 0 {
		query = query.Where("book_id = ?", req.BookID)
	}
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	var items []model.Item
	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(offset).Limit(req.Limit).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, module.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}
