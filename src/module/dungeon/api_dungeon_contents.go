package dungeon

import (
	"net/http"

	"github.com/bagaking/memorianexus/src/module/dto"

	"github.com/bagaking/goulp/wlog"
	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/gin-gonic/gin"
)

type ReqRemoveDungeonBooks struct {
	Books []utils.UInt64 `json:"books"`
}

type ReqRemoveDungeonItems struct {
	Items []utils.UInt64 `json:"items"`
}

type ReqRemoveDungeonTags struct {
	Tags []utils.UInt64 `json:"tags"`
}

// SubtractDungeonBooks handles removing books from a specific dungeon
// @Summary Remove books from a specific dungeon
// @Description 删除复习计划的 Books
// @TagNames dungeon
// @Accept json
// @Produce json
// @Param id path uint64 true "Dungeon ID"
// @Param books body ReqRemoveDungeonBooks true "Dungeon books data"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse "Invalid request body"
// @Failure 404 {object} utils.ErrorResponse "Dungeon or Book not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /dungeon/dungeons/{id}/books [delete]
func (svr *Service) SubtractDungeonBooks(c *gin.Context) {
	var req ReqRemoveDungeonBooks
	var dungeon model.Dungeon
	id := c.Param("id")

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinHandleError(c, wlog.ByCtx(c), http.StatusBadRequest, err, "Invalid request body")
		return
	}

	if err := svr.db.Where("id = ?", id).First(&dungeon).Error; err != nil {
		utils.GinHandleError(c, wlog.ByCtx(c), http.StatusNotFound, err, "Dungeon not found")
		return
	}

	for _, bookID := range req.Books {
		// 删除关联
		if err := svr.db.Where("dungeon_id = ? AND book_id = ?", dungeon.ID, bookID).Delete(&model.DungeonBook{}).Error; err != nil {
			utils.GinHandleError(c, wlog.ByCtx(c), http.StatusInternalServerError, err, "Failed to remove book from dungeon")
			return
		}
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Books removed from dungeon",
	})
}

// GetDungeonBooksDetail handles fetching the books of a specific dungeon
// @Summary Get the books of a specific dungeon
// @Description 获取复习计划的 Books
// @TagNames dungeon
// @Produce json
// @Param id path uint64 true "Dungeon ID"
// @Success 200 {array} utils.UInt64
// @Failure 404 {object} utils.ErrorResponse "Dungeon not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /dungeon/dungeons/{id}/books [get]
func (svr *Service) GetDungeonBooksDetail(c *gin.Context) {
	var dungeon model.Dungeon
	id := c.Param("id")

	if err := svr.db.Where("id = ?", id).First(&dungeon).Error; err != nil {
		utils.GinHandleError(c, wlog.ByCtx(c), http.StatusNotFound, err, "Dungeon not found")
		return
	}

	books, err := model.GetDungeonBookIDs(svr.db, dungeon.ID)
	if err != nil {
		utils.GinHandleError(c, wlog.ByCtx(c), http.StatusInternalServerError, err, "Failed to fetch dungeon books")
		return
	}

	c.JSON(http.StatusOK, books)
}

// SubtractDungeonTags handles removing tags from a specific dungeon
// @Summary Remove tags from a specific dungeon
// @Description 删除复习计划的 TagNames
// @TagNames dungeon
// @Accept json
// @Produce json
// @Param id path uint64 true "Dungeon ID"
// @Param tags body ReqRemoveDungeonTags true "Dungeon tags data"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse "Invalid request body"
// @Failure 404 {object} utils.ErrorResponse "Dungeon or Tag not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /dungeon/dungeons/{id}/tags [delete]
func (svr *Service) SubtractDungeonTags(c *gin.Context) {
	var req ReqRemoveDungeonTags
	var dungeon model.Dungeon
	id := c.Param("id")

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinHandleError(c, wlog.ByCtx(c), http.StatusBadRequest, err, "Invalid request body")
		return
	}

	if err := svr.db.Where("id = ?", id).First(&dungeon).Error; err != nil {
		utils.GinHandleError(c, wlog.ByCtx(c), http.StatusNotFound, err, "Dungeon not found")
		return
	}

	for _, tagID := range req.Tags {
		// 删除关联
		if err := svr.db.Where("dungeon_id = ? AND tag_id = ?", dungeon.ID, tagID).Delete(&model.DungeonTag{}).Error; err != nil {
			utils.GinHandleError(c, wlog.ByCtx(c), http.StatusInternalServerError, err, "Failed to remove tag from dungeon")
			return
		}
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "TagNames removed from dungeon",
	})
}

// GetDungeonTagsDetail handles fetching the tags of a specific dungeon
// @Summary Get the tags of a specific dungeon
// @Description 获取复习计划的 TagNames
// @TagNames dungeon
// @Produce json
// @Param id path uint64 true "Dungeon ID"
// @Success 200 {array} utils.UInt64
// @Failure 404 {object} utils.ErrorResponse "Dungeon not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /dungeon/dungeons/{id}/tags [get]
func (svr *Service) GetDungeonTagsDetail(c *gin.Context) {
	var dungeon model.Dungeon
	id := c.Param("id")

	if err := svr.db.Where("id = ?", id).First(&dungeon).Error; err != nil {
		utils.GinHandleError(c, wlog.ByCtx(c), http.StatusNotFound, err, "Dungeon not found")
		return
	}

	tags, err := model.GetDungeonTagIDs(svr.db, dungeon.ID)
	if err != nil {
		utils.GinHandleError(c, wlog.ByCtx(c), http.StatusInternalServerError, err, "Failed to fetch dungeon tags")
		return
	}

	c.JSON(http.StatusOK, tags)
}

// GetDungeonItemsDetail handles fetching the items of a specific dungeon
// @Summary Get the items of a specific dungeon
// @Description 获取复习计划的 Items
// @TagNames dungeon
// @Produce json
// @Param id path uint64 true "Dungeon ID"
// @Success 200 {array} utils.UInt64
// @Failure 404 {object} utils.ErrorResponse "Dungeon not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /dungeon/dungeons/{id}/items [get]
func (svr *Service) GetDungeonItemsDetail(c *gin.Context) {
	var dungeon model.Dungeon
	id := c.Param("id")

	if err := svr.db.Where("id = ?", id).First(&dungeon).Error; err != nil {
		utils.GinHandleError(c, wlog.ByCtx(c), http.StatusNotFound, err, "Dungeon not found")
		return
	}

	items, err := model.GetDungeonItemIDs(svr.db, dungeon.ID)
	if err != nil {
		utils.GinHandleError(c, wlog.ByCtx(c), http.StatusInternalServerError, err, "Failed to fetch dungeon tags")
		return
	}

	c.JSON(http.StatusOK, items)
}

// SubtractDungeonItems handles removing items from a specific dungeon
// @Summary Remove items from a specific dungeon
// @Description 删除复习计划的 Items
// @TagNames dungeon
// @Accept json
// @Produce json
// @Param id path uint64 true "Dungeon ID"
// @Param items body ReqRemoveDungeonItems true "Dungeon items data"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse "Invalid request body"
// @Failure 404 {object} utils.ErrorResponse "Dungeon or Item not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /dungeon/dungeons/{id}/items [delete]
func (svr *Service) SubtractDungeonItems(c *gin.Context) {
	var req ReqRemoveDungeonItems
	var dungeon model.Dungeon
	id := c.Param("id")

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinHandleError(c, wlog.ByCtx(c), http.StatusBadRequest, err, "Invalid request body")
		return
	}

	if err := svr.db.Where("id = ?", id).First(&dungeon).Error; err != nil {
		utils.GinHandleError(c, wlog.ByCtx(c), http.StatusNotFound, err, "Dungeon not found")
		return
	}

	for _, itemID := range req.Items {
		// 删除关联
		if err := svr.db.Where("dungeon_id = ? AND item_id = ?", dungeon.ID, itemID).Delete(&model.DungeonMonster{}).Error; err != nil {
			utils.GinHandleError(c, wlog.ByCtx(c), http.StatusInternalServerError, err, "Failed to remove item from dungeon")
			return
		}
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Items removed from dungeon",
	})
}
