package item

import (
	"net/http"

	"github.com/bagaking/goulp/wlog"
	"github.com/bagaking/memorianexus/internal/util"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/bagaking/memorianexus/src/module"
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
// @Tags item
// @Accept json
// @Produce json
// @Param id path uint64 true "Item ID"
// @Param item body ReqUpdateItem true "Item update data"
// @Success 200 {object} module.SuccessResponse "Successfully updated item"
// @Failure 400 {object} module.ErrorResponse "Bad Request with invalid item ID or update data"
// @Failure 500 {object} module.ErrorResponse "Internal Server Error with failing to update the item"
// @Router /items/{id} [put]
func (svr *Service) UpdateItem(c *gin.Context) {
	log := wlog.ByCtx(c, "UpdateItem")
	userID, exists := util.GetUIDFromGinCtx(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id, err := util.ParseIDFromString(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, module.ErrorResponse{Message: "Invalid item ID"})
		return
	}
	var req ReqUpdateItem
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, module.ErrorResponse{Message: err.Error()})
		return
	}

	// 开始数据库事务
	tx := svr.db.Begin()

	// 更新 Item 基础属性
	if err = tx.Model(&model.Item{}).
		Where("id = ?", id).
		Where("user_id = ?", userID).
		Updates(model.Item{
			Type:    req.Type,
			Content: req.Content,
		}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, module.ErrorResponse{Message: err.Error()})
		return
	}

	if err = updateItemTagsRef(c, tx, id, req.Tags); err != nil {
		log.WithError(err).Error("update item tags failed")
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, module.ErrorResponse{Message: err.Error()})
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, module.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, module.SuccessResponse{Message: "Item updated"})
}
