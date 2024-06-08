package item

import (
	"net/http"

	"github.com/bagaking/memorianexus/src/def"

	"github.com/bagaking/memorianexus/internal/utils"

	"github.com/bagaking/memorianexus/src/module/dto"

	"github.com/bagaking/goulp/wlog"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/gin-gonic/gin"
)

type ReqUpdateItem struct {
	Type       string              `json:"type,omitempty"`
	Content    string              `json:"content,omitempty"`
	Difficulty def.DifficultyLevel `json:"difficulty,omitempty"` // 难度，默认值为 NoviceNormal (0x01)
	Importance def.ImportanceLevel `json:"importance,omitempty"` // 重要程度，默认值为 DomainGeneral (0x01)
	Tags       []string            `json:"tags,omitempty"`       // 新增字段
}

// UpdateItem handles updating an existing item's information and associated tags.
// @Summary Update an item
// @Description Update an item's type, content, or associated tags.
// @TagNames item
// @Accept json
// @Produce json
// @Param id path uint64 true "Item ID"
// @Param item body ReqUpdateItem true "Item update data"
// @Success 200 {object} dto.RespItemUpdate "the updater"
// @Failure 400 {object} utils.ErrorResponse "Bad Request with invalid item ID or update data"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error with failing to update the item"
// @Router /items/{id} [put]
func (svr *Service) UpdateItem(c *gin.Context) {
	userID := utils.GinMustGetUserID(c)
	id := utils.GinMustGetID(c)
	log := wlog.ByCtx(c, "UpdateItem").WithField("user_id", userID).WithField("item_id", id)

	var req ReqUpdateItem
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, err, "invalid request data")
		return
	}

	updater := &model.Item{
		Type:       req.Type,
		Content:    req.Content,
		Difficulty: req.Difficulty,
		Importance: req.Importance,
	}
	// todo: 懒求值 update dungeon-monster 宽表冗余

	// 开始数据库事务
	tx := svr.db.Begin()

	if err := tx.Model(updater).Where("creator_id = ? AND id = ?", userID, id).Updates(updater).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "failed to update item")
		tx.Rollback()
		return
	}

	// 更新 Item 的 tags
	if err := model.UpdateItemTagsRef(c, tx, id, req.Tags); err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "failed to update item tags")
		tx.Rollback()
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "failed to commit transaction")
		tx.Rollback()
		return
	}

	new(dto.RespItemUpdate).With(new(dto.Item).FromModel(updater)).Response(c, "item updated")
}
