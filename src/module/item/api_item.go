package item

import (
	"errors"
	"net/http"

	"gorm.io/gorm"

	"github.com/bagaking/goulp/wlog"
	"github.com/gin-gonic/gin"
	"github.com/khicago/irr"

	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/bagaking/memorianexus/src/module/dto"
)

// GetItem handles retrieving a single item by ID, including its tags.
// @Summary Get an item by ID
// @Description Get detailed information about an item, including its tags.
// @Tags item
// @Accept json
// @Produce json
// @Param id path uint64 true "Item ID"
// @Success 200 {object} dto.RespItemGet "Successfully retrieved item with tags"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Router /items/{id} [get]
func (svr *Service) GetItem(c *gin.Context) {
	userID := utils.GinMustGetUserID(c)
	id := utils.GinMustGetID(c)
	log := wlog.ByCtx(c, "DeleteItem").WithField("user_id", userID).WithField("item_id", id)

	var item model.Item
	if err := svr.db.First(&item, id).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Cannot find item")
		return
	}

	// 获取 item 相关的 tags
	tags, err := model.GetItemTagNames(c, svr.db, item.ID)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Get item tag names failed")
		return
	}

	new(dto.RespItemGet).With(new(dto.Item).FromModel(&item, tags...)).Response(c, "item found")
}

// DeleteItem handles the deletion of an item.
// @Summary Delete an item
// @Description Delete an item from the system by ID.
// @TagNames item
// @Accept json
// @Produce json
// @Param id path uint64 true "Item ID"
// @Success 200 {object} dto.RespItemDelete "Successfully deleted item"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Router /items/{id} [delete]
func (svr *Service) DeleteItem(c *gin.Context) {
	userID := utils.GinMustGetUserID(c)
	id := utils.GinMustGetID(c)
	log := wlog.ByCtx(c, "DeleteItem").WithField("user_id", userID).WithField("item_id", id)

	// 在删除之前验证Item是否存在，并且属于当前操作的用户
	var item model.Item

	if err := svr.db.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.GinHandleError(c, log, http.StatusNotFound, err, "item not found")
		} else {
			utils.GinHandleError(c, log, http.StatusInternalServerError, err, "fetch item failed")
		}
		return
	}

	// 如果Item的UserID与请求中的用户ID不匹配，则拒绝删除操作
	if item.CreatorID != userID {
		err := irr.Error("unauthorized to delete this item, creator= %d", item.CreatorID)
		utils.GinHandleError(c, log, http.StatusForbidden, err, "unauthorized to delete this item")
		return
	}

	// 执行删除操作
	if err := svr.db.Delete(&model.Item{}, id).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "delete item failed")
		return
	}

	// 创建 DTO 并返回
	new(dto.RespItemDelete).With((&dto.Item{}).FromModel(&item)).Response(c, "item deleted")
}
