package tag

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bagaking/goulp/wlog"
	"github.com/khicago/got/util/typer"

	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/bagaking/memorianexus/src/module/dto"
)

// GetTags handles retrieving a list of all tags.
// @Summary Get all tags
// @Description Retrieves a list of all tags.
// @Tags tag
// @Accept json
// @Produce json
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Number of items per page"
// @Success 200 {array} dto.RespTagList "Successfully retrieved tags"
// @Router /tags [get]
func (svr *Service) GetTags(c *gin.Context) {
	pager := utils.GinGetPagerFromQuery(c)
	userID := utils.GinMustGetUserID(c)
	log := wlog.ByCtx(c, "GetTags").WithField("pager", pager).WithField("user_id", userID)

	tags, err := model.TagModel().GetTagsByUser(c, userID)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "failed to fetch tags")
		return
	}

	new(dto.RespTagList).WithPager(pager).Append(tags...).Response(c, "found tags")
}

// GetBooksByTag handles retrieving a list of books associated with a specific tag name.
// @Summary Get books by tag name
// @Description Retrieves a list of books associated with a specific tag name.
// @Tags tag
// @Accept json
// @Produce json
// @Param name path string true "Tag Name"
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Number of items per page"
// @Success 200 {array} dto.RespBookList "Successfully retrieved books"
// @Router /tags/{tag}/books [get]
func (svr *Service) GetBooksByTag(c *gin.Context) {
	userID := utils.GinMustGetUserID(c)
	tag := utils.GinMustGetTAG(c)
	pager := utils.GinGetPagerFromQuery(c)

	log := wlog.ByCtx(c, "GetBooksByTag").
		WithField("pager", pager).WithField("tag", tag).WithField("user_id", userID)

	books, err := model.FindBooksOfTag(c, svr.db, userID, tag, pager)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to fetch items by tag")
		return
	}

	new(dto.RespBookList).WithPager(pager).Append(typer.SliceMap(books, func(from model.Book) *dto.Book {
		return new(dto.Book).FromModel(&from)
	})...).Response(c, "found books by tag name")
}

// GetItemsByTag handles retrieving a list of items associated with a specific tag.
// @Summary Get items by tag
// @Description Retrieves a list of items associated with a specific tag.
// @Tags tag
// @Accept json
// @Produce json
// @Param id path uint64 true "Tag ID"
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Number of items per page"
// @Success 200 {array} dto.RespItemList "Successfully retrieved items"
// @Router /tags/{id}/items [get]
func (svr *Service) GetItemsByTag(c *gin.Context) {
	userID := utils.GinMustGetUserID(c)
	tag := utils.GinMustGetTAG(c)
	pager := utils.GinGetPagerFromQuery(c)

	log := wlog.ByCtx(c, "GetItemsByTag").
		WithField("pager", pager).WithField("tag", tag).WithField("user_id", userID)

	items, err := model.FindItemsOfTag(c, svr.db, userID, tag, pager)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to fetch items by tag")
		return
	}

	new(dto.RespItemList).WithPager(pager).Append(typer.SliceMap(items, func(from model.Item) *dto.Item {
		return new(dto.Item).FromModel(&from)
	})...).Response(c, "found items by tag")
}
