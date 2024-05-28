package item

import (
	"errors"
	"net/http"

	"github.com/bagaking/memorianexus/src/module/dto"

	"github.com/bagaking/goulp/wlog"
	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/gin-gonic/gin"
)

type ReqUpdateItem struct {
	Type    string   `json:"type,omitempty"`
	Content string   `json:"content,omitempty"`
	Tags    []string `json:"tags,omitempty"` // 新增字段
}

// UpdateItem handles updating an existing item's information and associated tags.
// @Summary Update an item
// @Description Update an item's type, content, or associated tags.
// @TagNames item
// @Accept json
// @Produce json
// @Param id path uint64 true "Item ID"
// @Param item body ReqUpdateItem true "Item update data"
// @Success 200 {object} dto.RespItemUpdate "Successfully updated item"
// @Failure 400 {object} utils.ErrorResponse "Bad Request with invalid item ID or update data"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error with failing to update the item"
// @Router /items/{id} [put]
func (svr *Service) UpdateItem(c *gin.Context) {
	log := wlog.ByCtx(c, "UpdateItem")
	userID, exists := utils.GetUIDFromGinCtx(c)
	if !exists {
		utils.GinHandleError(c, log, http.StatusUnauthorized, errors.New("user not authenticated"), "User not authenticated")
		return
	}

	id, err := utils.ParseIDFromString(c.Param("id"))
	if err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, err, "Invalid item ID")
		return
	}
	var req ReqUpdateItem
	if err = c.ShouldBindJSON(&req); err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, err, "Invalid request data")
		return
	}

	// 开始数据库事务
	tx := svr.db.Begin()

	// 更新 Item 基础属性
	if err = tx.Model(&model.Item{}).
		Where("id = ?", id).
		Where("creator_id = ?", userID).
		Updates(model.Item{
			Type:    req.Type,
			Content: req.Content,
		}).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to update item")
		tx.Rollback()
		return
	}

	// 更新 Item 的 tags
	if err = model.UpdateItemTagsRef(c, tx, id, req.Tags); err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to update item tags")
		tx.Rollback()
		return
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to commit transaction")
		tx.Rollback()
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, dto.RespItemUpdate{
		Message: "item updated",
	})
}
