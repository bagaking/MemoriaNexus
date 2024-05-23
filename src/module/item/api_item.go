package item

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/bagaking/goulp/wlog"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/bagaking/memorianexus/src/module"
)

// ItemDTO 数据传输对象
type ItemDTO struct {
	ID        utils.UInt64 `json:"id"`
	UserID    utils.UInt64 `json:"user_id"`
	Type      string       `json:"type"`
	Content   string       `json:"content"`
	Tags      []string     `json:"tags,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// GetItem handles retrieving a single item by ID, including its tags.
// @Summary Get an item by ID
// @Description Get detailed information about an item, including its tags.
// @Tags item
// @Accept json
// @Produce json
// @Param id path uint64 true "Item ID"
// @Success 200 {object} ItemDTO "Successfully retrieved item with tags"
// @Failure 400 {object} module.ErrorResponse "Bad Request"
// @Router /items/{id} [get]
func (svr *Service) GetItem(c *gin.Context) {
	log := wlog.ByCtx(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, err, "Invalid item ID")
		return
	}
	var item model.Item
	if err = svr.db.First(&item, id).Error; err != nil {
		utils.GinHandleError(c, log.WithField("item_id", id),
			http.StatusInternalServerError, err, "Cannot find item")
		return
	}

	// 获取 item 相关的 tags
	tags, err := model.GetItemTagNames(c, svr.db, item.ID)
	if err != nil {
		utils.GinHandleError(c, log.WithField("item_id", id),
			http.StatusInternalServerError, err, "Get item tag names failed")
		return
	}

	// 创建 DTO 并返回
	response := ItemDTO{
		ID:        item.ID,
		UserID:    item.UserID,
		Type:      item.Type,
		Content:   item.Content,
		Tags:      tags,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteItem handles the deletion of an item.
// @Summary Delete an item
// @Description Delete an item from the system by ID.
// @Tags item
// @Accept json
// @Produce json
// @Param id path uint64 true "Item ID"
// @Success 200 {object} module.SuccessResponse "Successfully deleted item"
// @Failure 400 {object} module.ErrorResponse "Bad Request"
// @Router /items/{id} [delete]
func (svr *Service) DeleteItem(c *gin.Context) {
	log := wlog.ByCtx(c)
	// 从请求上下文中获取当前用户ID
	userID, exists := utils.GetUIDFromGinCtx(c)
	if !exists {
		utils.GinHandleError(c, log, http.StatusUnauthorized, errors.New("user not authenticated"), "User not authenticated")
		return
	}

	// 解析URL中的Item ID
	itemID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.GinHandleError(c, log,
			http.StatusBadRequest, err, "Invalid item ID")
		return
	}

	// 在删除之前验证Item是否存在，并且属于当前操作的用户
	var item model.Item
	result := svr.db.First(&item, itemID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			utils.GinHandleError(c, log.WithField("item_id", itemID),
				http.StatusNotFound, result.Error, "Item not found")
		} else {
			utils.GinHandleError(c, log.WithField("item_id", itemID),
				http.StatusInternalServerError, result.Error, "Cannot find item")
		}
		return
	}

	// 如果Item的UserID与请求中的用户ID不匹配，则拒绝删除操作
	if item.UserID != userID {
		utils.GinHandleError(c, log.WithField("user_id", userID),
			http.StatusForbidden, errors.New("unauthorized to delete this item"), "Unauthorized to delete this item")
		return
	}

	// 执行删除操作
	if err = svr.db.Delete(&model.Item{}, itemID).Error; err != nil {
		utils.GinHandleError(c, log.WithField("item_id", itemID),
			http.StatusInternalServerError, err, "Failed to delete item")
		return
	}

	dto := ItemDTO{
		ID:        item.ID,
		UserID:    item.UserID,
		Type:      item.Type,
		Content:   item.Content,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}

	// 返回成功响应
	c.JSON(http.StatusOK, module.SuccessResponse{Message: "Item deleted", Data: dto})
}
