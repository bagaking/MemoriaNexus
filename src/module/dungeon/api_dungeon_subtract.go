package dungeon

import (
	"errors"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"github.com/bagaking/goulp/wlog"
	"github.com/gin-gonic/gin"
	"github.com/khicago/irr"

	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/bagaking/memorianexus/src/module/dto"
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
// @Tags dungeon
// @Accept json
// @Produce json
// @Param id path string true "Dungeon ID"
// @Param books body ReqRemoveDungeonBooks true "Dungeon books data"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse "Invalid request body"
// @Failure 404 {object} utils.ErrorResponse "Dungeon or Book not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /dungeon/dungeons/{id}/books [delete]
func (svr *Service) SubtractDungeonBooks(c *gin.Context) {
	log := wlog.ByCtx(c, "SubtractDungeonBooks")
	dungeonID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, err, "Invalid dungeon ID")
		return
	}

	var req ReqRemoveDungeonBooks
	if err = c.ShouldBindJSON(&req); err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, irr.Wrap(err, "parse request body failed"), "Invalid request body")
		return
	}

	dungeon := model.Dungeon{ID: utils.UInt64(dungeonID)}
	if err = svr.db.First(&dungeon).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.GinHandleError(c, log, http.StatusNotFound, err, "Dungeon not found")
		} else {
			utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to fetch dungeon")
		}
		return
	}

	for _, bookID := range req.Books {
		// 删除关联
		if err = svr.db.Where("dungeon_id = ? AND book_id = ?", dungeon.ID, bookID).Delete(&model.DungeonBook{}).Error; err != nil {
			utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to remove book from dungeon")
			return
		}
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Books removed from dungeon",
	})
}

// SubtractDungeonTags handles removing tags from a specific dungeon
// @Summary Remove tags from a specific dungeon
// @Description 删除复习计划的 TagNames
// @Tags dungeon
// @Accept json
// @Produce json
// @Param id path string true "Dungeon ID"
// @Param tags body ReqRemoveDungeonTags true "Dungeon tags data"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse "Invalid request body"
// @Failure 404 {object} utils.ErrorResponse "Dungeon or Tag not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /dungeon/dungeons/{id}/tags [delete]
func (svr *Service) SubtractDungeonTags(c *gin.Context) {
	log := wlog.ByCtx(c, "SubtractDungeonTags")
	dungeonID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, err, "Invalid dungeon ID")
		return
	}

	var req ReqRemoveDungeonTags
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, irr.Wrap(err, "parse request body failed"), "Invalid request body")
		return
	}

	dungeon := model.Dungeon{ID: utils.UInt64(dungeonID)}
	if err = svr.db.First(&dungeon).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.GinHandleError(c, log, http.StatusNotFound, err, "Dungeon not found")
		} else {
			utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to fetch dungeon")
		}
		return
	}

	for _, tagID := range req.Tags {
		// 删除关联
		if err = svr.db.Where("dungeon_id = ? AND tag_id = ?", dungeon.ID, tagID).Delete(&model.DungeonTag{}).Error; err != nil {
			utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to remove tag from dungeon")
			return
		}
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Tags removed from dungeon",
	})
}

// SubtractDungeonItems handles removing items from a specific dungeon
// @Summary Remove items from a specific dungeon
// @Description 删除复习计划的 Items
// @Tags dungeon
// @Accept json
// @Produce json
// @Param id path string true "Dungeon ID"
// @Param items body ReqRemoveDungeonItems true "Dungeon items data"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse "Invalid request body"
// @Failure 404 {object} utils.ErrorResponse "Dungeon or Item not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /dungeon/dungeons/{id}/items [delete]
func (svr *Service) SubtractDungeonItems(c *gin.Context) {
	log := wlog.ByCtx(c, "SubtractDungeonItems")
	dungeonID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, err, "Invalid dungeon ID")
		return
	}

	var req ReqRemoveDungeonItems
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, irr.Wrap(err, "parse request body failed"), "Invalid request body")
		return
	}

	dungeon := model.Dungeon{ID: utils.UInt64(dungeonID)}
	if err = svr.db.First(&dungeon).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.GinHandleError(c, log, http.StatusNotFound, err, "Dungeon not found")
		} else {
			utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to fetch dungeon")
		}
		return
	}

	for _, itemID := range req.Items {
		// 删除关联
		if err = svr.db.Where("dungeon_id = ? AND item_id = ?", dungeon.ID, itemID).Delete(&model.DungeonMonster{}).Error; err != nil {
			utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to remove item from dungeon")
			return
		}
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Items removed from dungeon",
	})
}
