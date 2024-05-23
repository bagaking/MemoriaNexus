package item

import (
	"errors"
	"net/http"

	"github.com/bagaking/goulp/wlog"

	"github.com/gin-gonic/gin"

	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/model"
)

// ReqGetItems encapsulates the request parameters for fetching items.
type ReqGetItems struct {
	UserID utils.UInt64 `form:"user_id"`
	BookID utils.UInt64 `form:"book_id"`
	Type   string       `form:"type"`
	Page   int          `form:"page"`
	Limit  int          `form:"limit"`
}

type RespItems struct {
	Items []ItemDTO `json:"items"`
	Page  int       `json:"page"`
	Limit int       `json:"limit"`
	Total int64     `json:"total"`
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
// @Success 200 {object} RespItems "Successfully retrieved items"
// @Failure 400 {object} module.ErrorResponse "Bad Request"
// @Router /items [get]
func (svr *Service) GetItems(c *gin.Context) {
	log := wlog.ByCtx(c, "GetItems")
	userID, exists := utils.GetUIDFromGinCtx(c)
	if !exists {
		utils.GinHandleError(c, log, http.StatusUnauthorized, errors.New("user not authenticated"), "User not authenticated")
		return
	}

	var req ReqGetItems
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

	query := svr.db.Model(&model.Item{})
	if req.UserID <= 0 { // 如果不指定用户，搜索的就是自己的
		req.UserID = userID
	}
	query = query.Where("user_id = ?", req.UserID)
	if req.BookID > 0 {
		query = query.Where("book_id = ?", req.BookID)
	}
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to count items")
		return
	}

	var items []model.Item
	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(offset).Limit(req.Limit).Find(&items).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to retrieve items")
		return
	}

	// 转换 Item 为 ItemDTO
	var itemDTOs []ItemDTO
	for _, item := range items {
		tags, err := model.GetItemTagNames(c, svr.db, item.ID)
		if err != nil {
			utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to get item tag names")
			return
		}
		itemDTO := ItemDTO{
			ID:        item.ID,
			UserID:    item.UserID,
			Type:      item.Type,
			Content:   item.Content,
			Tags:      tags,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
		itemDTOs = append(itemDTOs, itemDTO)
	}

	response := RespItems{
		Items: itemDTOs,
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
	}
	c.JSON(http.StatusOK, response)
}
