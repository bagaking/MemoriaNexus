package book

import (
	"github.com/bagaking/goulp/wlog"
	"github.com/gin-gonic/gin"
	"github.com/khicago/got/util/typer"
	"net/http"

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
	userID := utils.GinMustGetUserID(c)
	pager := utils.GinGetPagerFromQuery(c)
	log := wlog.ByCtx(c, "ListBooks").WithField("user_id", userID).WithField("pager", pager)

	var books []*model.Book
	result := svr.db.Where("user_id = ?", userID).Offset(pager.Offset).Limit(pager.Limit).Find(&books)

	if result.Error != nil {
		log.WithError(result.Error).Errorf("Failed to fetch books for user %v", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching books"})
		return
	}

	resp := new(dto.RespBooks).WithPager(pager).Append(
		typer.SliceMap(books, func(book *model.Book) dto.Book {
			return *(&dto.Book{}).FromModel(book)
		})...)
	resp.Response(c, "books found")
}
