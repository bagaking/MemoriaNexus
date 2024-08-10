package item

import (
	"errors"
	"net/http"

	"github.com/khicago/irr"
	"gorm.io/gorm"

	"github.com/bagaking/goulp/wlog"
	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/bagaking/memorianexus/src/module/dto"
	"github.com/gin-gonic/gin"
)

// CreateItem handles creating a new item with optional book affiliations and tags.
// @Summary Create a new item
// @Description Create a new item in the system with optional book affiliations and tags.
// @Tags item
// @Accept json
// @Produce json
// @Param item body ReqCreateItem true "Item creation data"
// @Success 201 {object} dto.RespItemCreate "Successfully created item with books and tags"
// @Failure 400 {object} utils.ErrorResponse "Bad Request if too many books or tags, or bad data"
// @Router /items [post]
func (svr *Service) CreateItem(c *gin.Context) {
	userID := utils.GinMustGetUserID(c)
	log := wlog.ByCtx(c, "CreateItem").WithField("user_id", userID)

	var req ReqCreateItem
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, err, "Invalid request data")
		return
	}

	// 检查 BookIDs 和 Tags 数量是否超出限制
	if len(req.BookIDs) > MaxBooksOncePerItem {
		utils.GinHandleError(c, log, http.StatusBadRequest, errors.New("too many books"), "Too many books")
		return
	}
	if len(req.Tags) > MaxTagsOncePerItem {
		utils.GinHandleError(c, log, http.StatusBadRequest, errors.New("too many tags"), "Too many tags")
		return
	}

	id, err := utils.GenIDU64(c)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to generate ID")
		return
	}

	// 创建 Item 实例
	item := &model.Item{
		ID:         id,
		CreatorID:  userID,
		Type:       req.Type,
		Content:    req.Content,
		Difficulty: req.Difficulty,
		Importance: req.Importance,
	}

	// 创建 Item 并开始数据库事务
	tx := svr.db.Begin()
	if err = tx.Create(item).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to create item")
		tx.Rollback()
		return
	}

	// 处理关联到 Book
	for _, bookID := range req.BookIDs {
		// 创建 Item 和 Book 的关系
		itemBook := &model.BookItem{
			ItemID: item.ID,
			BookID: bookID,
		}
		if err = svr.db.Create(itemBook).Error; err != nil {
			utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to create book link")
			tx.Rollback()
			return
		}
	}

	// 新增 Item 的 tags
	if req.Tags != nil && len(req.Tags) > 0 {
		if err = model.AddEntityTags(c, tx, userID, model.EntityTypeItem, id, req.Tags...); err != nil {
			utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to update item tags")
			tx.Rollback()
			return
		}
	}

	if err = tx.Commit().Error; err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to commit transaction")
		tx.Rollback()
		return
	}

	new(dto.RespItemCreate).With(new(dto.Item).FromModel(item, req.Tags...)).Response(c, "item created")
}

// ReadItem handles retrieving a single item by ID, including its tags.
// @Summary Get an item by ID
// @Description Get detailed information about an item, including its tags.
// @Tags item
// @Accept json
// @Produce json
// @Param id path uint64 true "Item ID"
// @Success 200 {object} dto.RespItemGet "Successfully retrieved item with tags"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Router /items/{id} [get]
func (svr *Service) ReadItem(c *gin.Context) {
	userID := utils.GinMustGetUserID(c)
	id := utils.GinMustGetID(c)
	log := wlog.ByCtx(c, "ReadItem").WithField("user_id", userID).WithField("item_id", id)

	var item model.Item
	if err := svr.db.First(&item, id).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Cannot find item")
		return
	}

	// 获取 item 相关的 tags
	tags, err := model.TagModel().GetTagsOfEntity(c, item.ID)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Get item tag names failed")
		return
	}

	new(dto.RespItemGet).With(new(dto.Item).FromModel(&item, tags...)).Response(c, "item found")
}

// UpdateItem handles updating an existing item's information and associated tags.
// @Summary Update an item
// @Description Update an item's type, content, or associated tags.
// @Tags item
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
	if err := model.UpdateEntityTagsDiff(c, tx, userID, id, req.Tags); err != nil {
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

// DeleteItem handles the deletion of an item.
// @Summary Delete an item
// @Description Delete an item from the system by ID.
// @Tags item
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
