package item

import (
	"errors"
	"net/http"

	"github.com/bagaking/goulp/wlog"
	"github.com/gin-gonic/gin"

	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/bagaking/memorianexus/src/module/dto"
)

type ReqCreateItem struct {
	Type    string         `json:"type,omitempty"`
	Content string         `json:"content,omitempty"`
	BookIDs []utils.UInt64 `json:"book_ids,omitempty"` // 用于接收一个或多个 BookID
	Tags    []string       `json:"tags,omitempty"`     // 新增字段，用于接收一组 Tag 名称
}

const (
	MaxBooksOncePerItem = 10 // 设定每个 Item 可以关联的最大 Books 数量
	MaxTagsOncePerItem  = 5  // 设定每个 Item 可以拥有的最大 TagNames 数量
)

// CreateItem handles creating a new item with optional book affiliations and tags.
// @Summary Create a new item
// @Description Create a new item in the system with optional book affiliations and tags.
// @TagNames item
// @Accept json
// @Produce json
// @Param item body ReqCreateItem true "Item creation data"
// @Success 201 {object} dto.RespItemCreate "Successfully created item with books and tags"
// @Failure 400 {object} utils.ErrorResponse "Bad Request if too many books or tags, or bad data"
// @Router /items [post]
func (svr *Service) CreateItem(c *gin.Context) {
	log := wlog.ByCtx(c, "CreateItem")
	userID, exists := utils.GetUIDFromGinCtx(c)
	if !exists {
		utils.GinHandleError(c, log, http.StatusUnauthorized, errors.New("user not authenticated"), "User not authenticated")
		return
	}

	var req ReqCreateItem
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, err, "Invalid request data")
		return
	}

	// 检查 BookIDs 和 TagNames 数量是否超出限制
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
		ID:        id,
		CreatorID: userID,
		Type:      req.Type,
		Content:   req.Content,
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
	c.JSON(http.StatusOK, dto.RespItemCreate{
		Message: "item created",
		Data:    (&dto.Item{}).FromModel(item, req.Tags...),
	})
}
