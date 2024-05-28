package dungeon

import (
	"net/http"

	"github.com/bagaking/goulp/wlog"
	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/bagaking/memorianexus/src/module/dto"
	"github.com/gin-gonic/gin"
)

type ReqAddDungeonBooks struct {
	Books []utils.UInt64 `json:"books"`
}

type ReqAddDungeonItems struct {
	Items []utils.UInt64 `json:"items"`
}

type ReqAddDungeonTags struct {
	Tags []utils.UInt64 `json:"tags"`
}

// AppendBooksToDungeon handles adding books to a specific dungeon
// @Summary Add books to a specific dungeon
// @Description 添加复习计划的 Books
// @TagNames dungeon
// @Accept json
// @Produce json
// @Param id path uint64 true "Dungeon ID"
// @Param books body ReqAddDungeonBooks true "Dungeon books data"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse "Invalid request body"
// @Failure 404 {object} utils.ErrorResponse "Dungeon or Book not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /dungeon/dungeons/{id}/books [post]
func (svr *Service) AppendBooksToDungeon(c *gin.Context) {
	log := wlog.ByCtx(c, "AppendTagsToDungeon")
	var req ReqAddDungeonBooks
	var dungeon model.Dungeon
	id := c.Param("id")

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	if err := svr.db.Where("id = ?", id).First(&dungeon).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusNotFound, err, "Dungeon not found")
		return
	}

	if err := dungeon.AddMonster(svr.db, model.MonsterSourceBook, req.Books); err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Books added to dungeon",
	})
}

// AppendTagsToDungeon handles adding tags to a specific dungeon
// @Summary Add tags to a specific dungeon
// @Description 添加复习计划的 TagNames
// @TagNames dungeon
// @Accept json
// @Produce json
// @Param id path uint64 true "Dungeon ID"
// @Param tags body ReqAddDungeonTags true "Dungeon tags data"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse "Invalid request body"
// @Failure 404 {object} utils.ErrorResponse "Dungeon or Tag not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /dungeon/dungeons/{id}/tags [post]
func (svr *Service) AppendTagsToDungeon(c *gin.Context) {
	log := wlog.ByCtx(c, "AppendTagsToDungeon")
	var req ReqAddDungeonTags
	var dungeon model.Dungeon
	id := c.Param("id")

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	if err := svr.db.Where("id = ?", id).First(&dungeon).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusNotFound, err, "Dungeon not found")
		return
	}

	if err := dungeon.AddMonster(svr.db, model.MonsterSourceBook, req.Tags); err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "TagNames added to dungeon",
	})
}

// AppendItemsToDungeon handles adding items to a specific dungeon
// @Summary Add items to a specific dungeon
// @Description 添加复习计划的 Items
// @TagNames dungeon
// @Accept json
// @Produce json
// @Param id path uint64 true "Dungeon ID"
// @Param items body ReqAddDungeonItems true "Dungeon items data"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse "Invalid request body"
// @Failure 404 {object} utils.ErrorResponse "Dungeon or Item not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /dungeon/dungeons/{id}/items [post]
func (svr *Service) AppendItemsToDungeon(c *gin.Context) {
	log := wlog.ByCtx(c, "AppendItemsToDungeon")
	var req ReqAddDungeonItems
	var dungeon model.Dungeon
	id := c.Param("id")

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	if err := svr.db.Where("id = ?", id).First(&dungeon).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusNotFound, err, "Dungeon not found")
		return
	}

	if err := dungeon.AddMonster(svr.db, model.MonsterSourceItem, req.Items); err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Items added to dungeon",
	})
}
