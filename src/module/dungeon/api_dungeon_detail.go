package dungeon

import (
	"net/http"

	"github.com/bagaking/goulp/wlog"
	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/gin-gonic/gin"
)

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
