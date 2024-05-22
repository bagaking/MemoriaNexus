package item

import (
	"net/http"

	"github.com/bagaking/goulp/wlog"
	"github.com/bagaking/memorianexus/internal/util"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/bagaking/memorianexus/src/module"
	"github.com/gin-gonic/gin"
)

type ReqCreateItem struct {
	Type    string        `json:"type,omitempty"`
	Content string        `json:"content,omitempty"`
	BookIDs []util.UInt64 `json:"book_ids,omitempty"` // 用于接收一个或多个 BookID
	Tags    []string      `json:"tags,omitempty"`     // 新增字段，用于接收一组 Tag 名称
}

const (
	MaxBooksOncePerItem = 10 // 设定每个 Item 可以关联的最大 Books 数量
	MaxTagsOncePerItem  = 5  // 设定每个 Item 可以拥有的最大 Tags 数量
)

// CreateItem handles creating a new item with optional book affiliations and tags.
// @Summary Create a new item
// @Description Create a new item in the system with optional book affiliations and tags.
// @Tags item
// @Accept json
// @Produce json
// @Param item body ReqCreateItem true "Item creation data"
// @Success 201 {object} model.Item "Successfully created item with books and tags"
// @Failure 400 {object} module.ErrorResponse "Bad Request if too many books or tags, or bad data"
// @Router /items [post]
func (svr *Service) CreateItem(c *gin.Context) {
	log := wlog.ByCtx(c, "CreateItem")
	userID, exists := util.GetUIDFromGinCtx(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req ReqCreateItem

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, module.ErrorResponse{Message: err.Error()})
		return
	}

	// 检查 BookIDs 和 Tags 数量是否超出限制
	if len(req.BookIDs) > MaxBooksOncePerItem {
		c.JSON(http.StatusBadRequest, module.ErrorResponse{Message: "Too many books"})
		return
	}
	if len(req.Tags) > MaxTagsOncePerItem {
		c.JSON(http.StatusBadRequest, module.ErrorResponse{Message: "Too many tags"})
		return
	}

	id, err := util.GenIDU64(c)
	if err != nil {
		log.WithError(err).Error("generate id failed")
		c.JSON(http.StatusInternalServerError, module.ErrorResponse{Message: "generate id failed"})
	}

	// 创建 Item 实例
	item := &model.Item{
		ID:      id,
		UserID:  userID,
		Type:    req.Type,
		Content: req.Content,
	}

	// 创建 Item 并开始数据库事务
	tx := svr.db.Begin()
	if err = tx.Create(item).Error; err != nil {
		tx.Rollback()
		log.WithError(err).Error("create item failed")
		c.JSON(http.StatusInternalServerError, module.ErrorResponse{Message: err.Error()})
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
			log.WithError(err).Error("create book link failed")
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, module.ErrorResponse{Message: err.Error()})
			return
		}
	}

	if err = updateItemTagsRef(c, tx, item.ID, req.Tags); err != nil {
		log.WithError(err).Error("update item tags failed")
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, module.ErrorResponse{Message: err.Error()})
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		log.WithError(err).Error("create item failed")
		c.JSON(http.StatusInternalServerError, module.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}