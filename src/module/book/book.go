package book

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/bagaking/goulp/wlog"
	"github.com/khicago/got/util/typer"
	"github.com/khicago/irr"

	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/bagaking/memorianexus/src/module/dto"
)

// CreateBook handles the request to create a new book.
// @Summary Create a book
// @Description Creates a new book and optionally associates tags with it
// @Tags book
// @Accept json
// @Produce json
// @Param book body dto.ReqCreateBook true "Book creation data"
// @Success 201 {object} dto.RespBookCreate "Successfully created book"
// @Failure 400 {object} utils.ErrorResponse "Invalid parameters"
// @Router /books [post]
func (svr *Service) CreateBook(c *gin.Context) {
	log := wlog.ByCtx(c, "CreateBook")
	userID := utils.GinMustGetUserID(c)

	var req dto.ReqCreateBook
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	bookID, err := utils.GenIDU64(c)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to generate ID")
		return
	}

	book := &model.Book{
		ID:          bookID,
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
	}

	// Begin a transaction
	tx := svr.db.Begin()
	defer func() {
		// Ensure transaction rollback on error or panic
		if r := recover(); r != nil || tx.Error != nil {
			tx.Rollback()
			if r != nil {
				utils.GinHandleError(c, log, http.StatusInternalServerError, fmt.Errorf("%v", r), "Transaction failed")
			} else {
				utils.GinHandleError(c, log, http.StatusInternalServerError, tx.Error, "Transaction failed")
			}
		}
	}()

	// Create the book record in the database
	if err = tx.Create(book).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			utils.GinHandleError(c, log, http.StatusConflict, err, "Book already exists")
		} else {
			utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to create book")
		}
		return
	}

	// Update tags associated with the book
	// todo: should not using one tx
	if err = model.UpdateBookTagsRef(c, tx, book.ID, req.Tags); err != nil {
		tx.Rollback()
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to update book tags")
		return
	}

	// Commit the transaction
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to commit transaction")
		return
	}

	// Construct the response
	resp := dto.RespBookCreate{
		Message: "book created",
		Data:    new(dto.Book).FromModel(book),
	}

	// Send the response
	c.JSON(http.StatusCreated, resp)
}

// GetBooks handles retrieving a list of books with pagination.
// @Summary Get list of books with pagination
// @Description Get a paginated list of books for the user.
// @Tags book
// @Accept json
// @Produce json
// @Param page query int false "Page number for pagination" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} dto.RespBooks "Successfully retrieved list of books"
// @Router /books [get]
func (svr *Service) GetBooks(c *gin.Context) {
	log := wlog.ByCtx(c, "GetBooks")
	userID := utils.GinMustGetUserID(c)

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		log.WithError(err).Error("Invalid page number")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.WithError(err).Error("Invalid limit number")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	var books []*model.Book
	result := svr.db.Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&books)

	if result.Error != nil {
		log.WithError(result.Error).Errorf("Failed to fetch books for user %v", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching books"})
		return
	}

	resp := new(dto.RespBooks).SetPageAndLimit(page, limit).Append(
		typer.SliceMap(books, func(book *model.Book) dto.Book {
			return *(&dto.Book{}).FromModel(book)
		})...)
	resp.Response(c, "books found")
}

// GetBook handles retrieving a single book by ID.
// @Summary Get a book by ID
// @Description Get detailed information about a book.
// @Tags book
// @Accept json
// @Produce json
// @Param id path uint64 true "Book ID"
// @Success 200 {object} dto.RespBookGet "Successfully retrieved book"
// @Router /books/{id} [get]
func (svr *Service) GetBook(c *gin.Context) {
	log := wlog.ByCtx(c, "GetBook")
	userID := utils.GinMustGetUserID(c)
	id := utils.GinMustGetID(c)

	book := &model.Book{}
	result := svr.db.Where("id = ? AND user_id = ?", id, userID).First(book)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found or permission denied"})
		} else {
			log.WithError(err).Errorf("Failed to fetch books for user %v", userID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching book"})
		}
		return
	}

	new(dto.RespBookGet).With(new(dto.Book).FromModel(book)).Response(c, "book found")
}

// UpdateBook handles updating a book's information.
// @Summary Update book information
// @Description Update information for an existing book.
// @Tags book
// @Accept json
// @Produce json
// @Param id path uint64 true "Book ID"
// @Param book body dto.ReqCreateBook true "Book update data"
// @Success 200 {object} dto.RespBookUpdate "Successfully updated book"
// @Router /books/{id} [put]
func (svr *Service) UpdateBook(c *gin.Context) {
	userID := utils.GinMustGetUserID(c)
	id := utils.GinMustGetID(c)

	log := wlog.ByCtx(c, "UpdateBook").WithField("user_id", userID).WithField("book_id", id)

	var req dto.ReqCreateBook
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, err, "invalid request body")
		return
	}

	// 尝试更新记录
	updater := &model.Book{
		Title:       req.Title,
		Description: req.Description,
	}
	result := svr.db.Model(updater).Where("id = ? AND user_id = ?", id, userID).Updates(updater)

	// 处理错误
	if err := result.Error; err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "updating book failed")
		return
	}

	// 检查是否实际修改了记录
	if result.RowsAffected == 0 {
		err := irr.Error("Book not found or permission denied")
		utils.GinHandleError(c, log, http.StatusNotFound, err, "nothing changed")
		return
	}

	new(dto.RespBookUpdate).With(new(dto.Book).FromModel(updater)).Response(c, "book updated")
}

// DeleteBook handles deleting a book.
// @Summary Delete a book
// @Description Delete a book from the system by ID.
// @Tags book
// @Accept json
// @Produce json
// @Param id path uint64 true "Book ID"
// @Success 200 {object} dto.RespBookDelete "Successfully deleted book"
// @Router /books/{id} [delete]
func (svr *Service) DeleteBook(c *gin.Context) {
	userID := utils.GinMustGetUserID(c)
	id := utils.GinMustGetID(c)

	log := wlog.ByCtx(c, "DeleteBook").WithField("user_id", userID).WithField("book_id", id)

	// 删除前校验所有者
	var book model.Book
	if err := svr.db.Where("id = ? AND user_id = ?", id, userID).First(&book).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusNotFound, err, "book not found or permission denied")
		return
	}

	// 开始事务
	tx := svr.db.Begin()

	// 删除书册与标签、项的关系
	if err := tx.Where("book_id = ?", id).Delete(&model.BookTag{}).Error; err != nil {
		tx.Rollback()
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "failed to delete book-tags")
		return
	}

	if err := tx.Where("book_id = ?", id).Delete(&model.BookItem{}).Error; err != nil {
		tx.Rollback()
		utils.GinHandleError(c, log, http.StatusInternalServerError,
			irr.Wrap(err, "user=%v book_id=%v", userID, id), "failed to delete book-items")
		return
	}

	// 删除书册
	if err := tx.Delete(&model.Book{}, id).Error; err != nil {
		tx.Rollback()
		utils.GinHandleError(c, log, http.StatusInternalServerError,
			irr.Wrap(err, "user=%v book_id=%v", userID, id), "failed to delete book")
		return
	}

	tx.Commit()

	new(dto.RespBookDelete).With(id).Response(c, "book deleted")
}
