package item

import (
	"errors"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

	"github.com/bagaking/memorianexus/internal/util"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/bagaking/memorianexus/src/module"
)

// GetItem handles retrieving a single item by ID.
// @Summary Get an item by ID
// @Description Get detailed information about an item.
// @Tags item
// @Accept json
// @Produce json
// @Param id path uint64 true "Item ID"
// @Success 200 {object} model.Item "Successfully retrieved item"
// @Failure 400 {object} module.ErrorResponse "Bad Request"
// @Router /items/{id} [get]
func (svr *Service) GetItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, module.ErrorResponse{Message: "Invalid item ID"})
		return
	}
	var item model.Item
	if err = svr.db.First(&item, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, module.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, &item)
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
	// 从请求上下文中获取当前用户ID
	userID, exists := util.GetUIDFromGinCtx(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 解析URL中的Item ID
	itemID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, module.ErrorResponse{Message: "Invalid item ID"})
		return
	}

	// 在删除之前验证Item是否存在，并且属于当前操作的用户
	var item model.Item
	result := svr.db.First(&item, itemID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, module.ErrorResponse{Message: "Item not found"})
		} else {
			c.JSON(http.StatusInternalServerError, module.ErrorResponse{Message: result.Error.Error()})
		}
		return
	}

	// 如果Item的UserID与请求中的用户ID不匹配，则拒绝删除操作
	if item.UserID != userID {
		c.JSON(http.StatusForbidden, module.ErrorResponse{Message: "Unauthorized to delete this item"})
		return
	}

	// 执行删除操作
	if err = svr.db.Delete(&model.Item{}, itemID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, module.ErrorResponse{Message: err.Error()})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, module.SuccessResponse{Message: "Item deleted"})
}
