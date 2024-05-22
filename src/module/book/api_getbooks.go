package book

import (
	"errors"
	"github.com/bagaking/goulp/wlog"
	"github.com/bagaking/memorianexus/internal/util"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/gin-gonic/gin"
	"github.com/khicago/got/util/typer"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// GetBooks handles retrieving a list of books with pagination.
// @Summary Get list of books with pagination
// @Description Get a paginated list of books for the user.
// @Tags book
// @Accept json
// @Produce json
// @Param page query int false "Page number for pagination" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} RespBooks "Successfully retrieved list of books"
// @Router /books [get]
func (svr *Service) GetBooks(c *gin.Context) {
	log := wlog.ByCtx(c, "GetBooks")
	userID, exists := util.GetUIDFromGinCtx(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

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

	var books []model.Book
	result := svr.db.Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&books)

	if result.Error != nil {
		log.WithError(result.Error).Errorf("Failed to fetch books for user %v", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching books"})
		return
	}

	resp := RespBooks{
		Books: typer.SliceMap(books, func(book model.Book) RespBook {
			resp := RespBook{
				ID:          book.ID,
				UserID:      book.UserID,
				Title:       book.Title,
				Description: book.Description,
				CreatedAt:   book.CreatedAt,
				UpdatedAt:   book.UpdatedAt,
			}
			return resp
		}),
		Page:  page,
		Limit: limit,
	}

	c.JSON(http.StatusOK, resp)
}

// GetBook handles retrieving a single book by ID.
// @Summary Get a book by ID
// @Description Get detailed information about a book.
// @Tags book
// @Accept json
// @Produce json
// @Param id path uint64 true "Book ID"
// @Success 200 {object} RespBook "Successfully retrieved book"
// @Router /books/{id} [get]
func (svr *Service) GetBook(c *gin.Context) {
	log := wlog.ByCtx(c, "GetBook")
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

	var book model.Book
	result := svr.db.Where("id = ? AND user_id = ?", id, userID).First(&book)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found or permission denied"})
		} else {
			log.WithError(err).Errorf("Failed to fetch books for user %v", userID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching book"})
		}
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

	c.JSON(http.StatusOK, resp)
}
