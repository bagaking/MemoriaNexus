package tag

import (
	"net/http"

	"github.com/bagaking/goulp/wlog"
	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/bagaking/memorianexus/src/module/dto"
	"github.com/gin-gonic/gin"
	"github.com/khicago/got/util/typer"
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
	log := wlog.ByCtx(c, "GetTags").WithField("pager", pager)

	var tags []model.Tag
	if err := svr.db.Find(&tags).Offset(pager.Offset).Limit(pager.Limit).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to fetch tags")
		return
	}

	new(dto.RespTagList).WithPager(pager).Append(typer.SliceMap(tags, func(from model.Tag) *dto.Tag {
		return new(dto.Tag).FromModel(from)
	})...).Response(c, "found tags")
}

// GetTagByName handles retrieving a tag by its name.
// @Summary Get tag by name
// @Description Retrieves a tag by its name.
// @Tags tag
// @Accept json
// @Produce json
// @Param name path string true "Tag Name"
// @Success 200 {object} dto.Tag "Successfully retrieved tag"
// @Router /tags/name/{name} [get]
func (svr *Service) GetTagByName(c *gin.Context) {
	tagName := c.Param("name")
	pager := utils.GinGetPagerFromQuery(c)
	log := wlog.ByCtx(c, "GetTagByName").WithField("pager", pager)

	var tag model.Tag
	if err := svr.db.Where("name = ?", tagName).First(&tag).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to fetch tag")
		return
	}

	new(dto.RespTagGet).With(new(dto.Tag).FromModel(tag)).Response(c, "found tag by name")
}

// GetBooksByTagName handles retrieving a list of books associated with a specific tag name.
// @Summary Get books by tag name
// @Description Retrieves a list of books associated with a specific tag name.
// @Tags tag
// @Accept json
// @Produce json
// @Param name path string true "Tag Name"
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Number of items per page"
// @Success 200 {array} dto.RespBookList "Successfully retrieved books"
// @Router /tags/name/{name}/books [get]
func (svr *Service) GetBooksByTagName(c *gin.Context) {
	log := wlog.ByCtx(c, "GetBooksByTagName")
	tagName := c.Param("name")
	pager := utils.GinGetPagerFromQuery(c)

	var tag model.Tag
	if err := svr.db.Where("name = ?", tagName).First(&tag).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to fetch tag")
		return
	}

	var bookTags []model.BookTag
	if err := svr.db.Where("tag_id = ?", tag.ID).Offset(pager.Offset).Limit(pager.Limit).Find(&bookTags).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to fetch book tags")
		return
	}

	var bookIDs []utils.UInt64
	for _, bookTag := range bookTags {
		bookIDs = append(bookIDs, bookTag.BookID)
	}

	var books []model.Book
	if err := svr.db.Where("id IN (?)", bookIDs).Find(&books).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to fetch books")
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
	log := wlog.ByCtx(c, "GetItemsByTag")
	tagID := utils.GinMustGetID(c)
	pager := utils.GinGetPagerFromQuery(c)

	var itemTags []model.ItemTag
	if err := svr.db.Where("tag_id = ?", tagID).Offset(pager.Offset).Limit(pager.Limit).Find(&itemTags).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to fetch item tags")
		return
	}

	var itemIDs []utils.UInt64
	for _, itemTag := range itemTags {
		itemIDs = append(itemIDs, itemTag.ItemID)
	}

	var items []model.Item
	if err := svr.db.Where("id IN (?)", itemIDs).Offset(pager.Offset).Limit(pager.Limit).Find(&items).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to fetch items")
		return
	}

	new(dto.RespItemList).WithPager(pager).Append(typer.SliceMap(items, func(from model.Item) *dto.Item {
		return new(dto.Item).FromModel(&from)
	})...).Response(c, "found items by tag")
}

// GetItemsByTagName handles retrieving a list of items associated with a specific tag name.
// @Summary Get items by tag name
// @Description Retrieves a list of items associated with a specific tag name.
// @Tags tag
// @Accept json
// @Produce json
// @Param name path string true "Tag Name"
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Number of items per page"
// @Success 200 {array} dto.RespItemList "Successfully retrieved items"
// @Router /tags/name/{name}/items [get]
func (svr *Service) GetItemsByTagName(c *gin.Context) {
	log := wlog.ByCtx(c, "GetItemsByTagName")
	tagName := c.Param("name")
	pager := utils.GinGetPagerFromQuery(c)

	var tag model.Tag
	if err := svr.db.Where("name = ?", tagName).First(&tag).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to fetch tag")
		return
	}

	var itemTags []model.ItemTag
	if err := svr.db.Where("tag_id = ?", tag.ID).Offset(pager.Offset).Limit(pager.Limit).Find(&itemTags).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to fetch item tags")
		return
	}

	var itemIDs []utils.UInt64
	for _, itemTag := range itemTags {
		itemIDs = append(itemIDs, itemTag.ItemID)
	}

	var items []model.Item
	if err := svr.db.Where("id IN (?)", itemIDs).Offset(pager.Offset).Limit(pager.Limit).Find(&items).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to fetch items")
		return
	}

	new(dto.RespItemList).WithPager(pager).Append(typer.SliceMap(items, func(from model.Item) *dto.Item {
		return new(dto.Item).FromModel(from)
	})...).Response(c, "found items by tag name")
}
