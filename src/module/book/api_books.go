package book

import (
	"net/http"
	"strconv"

	"github.com/bagaking/goulp/wlog"
	"github.com/gin-gonic/gin"
	"github.com/khicago/got/util/typer"

	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/bagaking/memorianexus/src/module/dto"
)

// ListBooks handles retrieving a list of books with pagination.
// @Summary Get list of books with pagination
// @Description Get a paginated list of books for the user.
// @Tags book
// @Accept json
// @Produce json
// @Param page query int false "Page number for pagination" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} dto.RespBooks "Successfully retrieved list of books"
// @Router /books [get]
func (svr *Service) ListBooks(c *gin.Context) {
	log := wlog.ByCtx(c, "ListBooks")
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
