package book

import (
	"github.com/bagaking/goulp/wlog"
	"github.com/bagaking/memorianexus/internal/util"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/gin-gonic/gin"
	"github.com/khicago/irr"
	"net/http"
	"strconv"
)

// UpdateBook handles updating a book's information.
// @Summary Update book information
// @Description Update information for an existing book.
// @Tags book
// @Accept json
// @Produce json
// @Param id path uint64 true "Book ID"
// @Param book body ReqCreateBook true "Book update data"
// @Success 200 {string} string "Successfully updated book"
// @Router /books/{id} [put]
func (svr *Service) UpdateBook(c *gin.Context) {
	log := wlog.ByCtx(c, "UpdateBook")
	userID, exists := util.GetUIDFromGinCtx(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var req ReqCreateBook
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 尝试更新记录
	result := svr.db.Model(&model.Book{}).Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]any{
			"title":       req.Title,
			"description": req.Description,
		})

	// 处理错误
	if err = result.Error; err != nil {
		log.WithError(err).Errorf("Failed to update books for user %v book_id= %v", userID, id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating book"})
		return
	}

	// 检查是否实际修改了记录
	if result.RowsAffected == 0 {
		err = irr.Error("Book not found or permission denied")
		log.WithError(err).Errorf("Failed to update books for user %v book_id= %v", userID, id)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book updated", "book_id": id})
}

// DeleteBook handles deleting a book.
// @Summary Delete a book
// @Description Delete a book from the system by ID.
// @Tags book
// @Accept json
// @Produce json
// @Param id path uint64 true "Book ID"
// @Success 200 {string} string "Successfully deleted book"
// @Router /books/{id} [delete]
func (svr *Service) DeleteBook(c *gin.Context) {
	log := wlog.ByCtx(c, "DeleteBook")
	userID, exists := util.GetUIDFromGinCtx(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	// 删除前校验所有者
	var book model.Book
	if err = svr.db.Where("id = ? AND user_id = ?", id, userID).First(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found or permission denied"})
		return
	}

	// 开始事务
	tx := svr.db.Begin()

	// 删除书册与标签、项的关系
	if err = tx.Where("book_id = ?", id).Delete(&model.BookTag{}).Error; err != nil {
		tx.Rollback()
		log.WithError(err).Errorf("Failed to delete book_tag_ref for user %v book_id= %v", userID, id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book tags"})
		return
	}

	if err = tx.Where("book_id = ?", id).Delete(&model.BookItem{}).Error; err != nil {
		tx.Rollback()
		log.WithError(err).Errorf("Failed to delete book_item_ref for user %v book_id= %v", userID, id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book items"})
		return
	}

	// 删除书册
	if err = tx.Delete(&model.Book{}, id).Error; err != nil {
		tx.Rollback()
		log.WithError(err).Errorf("Failed to delete book for user %v book_id= %v", userID, id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted", "book_id": id})
}
