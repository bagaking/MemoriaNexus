package book

import (
	"github.com/bagaking/goulp/wlog"
	"github.com/bagaking/memorianexus/internal/util"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/bagaking/memorianexus/src/module"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateBook 处理创建书册的请求
// @Summary 创建书册
// @Description 创建一个新的书册，可选地关联标签
// @Tags book
// @Accept json
// @Produce json
// @Param book body ReqCreateBook true "书册创建数据"
// @Success 201 {object} RespBook "成功创建书册"
// @Failure 400 {object} module.ErrorResponse "参数错误"
// @Router /books [post]
func (svr *Service) CreateBook(c *gin.Context) {
	log := wlog.ByCtx(c, "CreateBook")
	userID, exists := util.GetUIDFromGinCtx(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	var req ReqCreateBook
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, module.ErrorResponse{Message: err.Error()})
		return
	}

	bookID, err := util.GenIDU64(c)
	if err != nil {
		log.WithError(err).Error("生成ID失败")
		c.JSON(http.StatusInternalServerError, module.ErrorResponse{Message: "生成ID失败"})
		return
	}

	book := &model.Book{
		ID:          bookID,
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
	}

	// 开始数据库事务
	tx := svr.db.Begin()
	if err = tx.Create(book).Error; err != nil {
		tx.Rollback()
		log.WithError(err).Error("创建书册失败")
		c.JSON(http.StatusInternalServerError, module.ErrorResponse{Message: err.Error()})
		return
	}

	// 处理关联标签
	if err = updateBookTagsRef(c, tx, book.ID, req.Tags); err != nil {
		log.WithError(err).Error("更新书册标签失败")
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, module.ErrorResponse{Message: err.Error()})
		return
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, module.ErrorResponse{Message: err.Error()})
		return
	}

	resp := RespBook{
		ID:          book.ID,
		UserID:      book.UserID,
		Title:       book.Title,
		Description: book.Description,
		CreatedAt:   book.CreatedAt,
		UpdatedAt:   book.UpdatedAt,
	}

	c.JSON(http.StatusCreated, resp)
}
