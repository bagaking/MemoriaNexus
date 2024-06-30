package dungeon

import (
	"errors"
	"net/http"

	"gorm.io/gorm"

	"github.com/bagaking/memorianexus/src/module/dto"

	"github.com/gin-gonic/gin"

	"github.com/bagaking/goulp/wlog"
	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/model"
)

// GetDungeonBooksDetail handles fetching the books of a specific dungeon
// @Summary Get the books of a specific dungeon
// @Description 获取复习计划的 Books
// @Tags dungeon
// @Produce json
// @Param id path uint64 true "Dungeon ID"
// @Success 200 {array} utils.UInt64
// @Failure 404 {object} utils.ErrorResponse "Dungeon not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /dungeon/dungeons/{id}/books [get]
func (svr *Service) GetDungeonBooksDetail(c *gin.Context) {
	userID, dungeonID := utils.GinMustGetUserID(c), utils.GinMustGetID(c)
	log := wlog.ByCtx(c, "GetDungeonBooksDetail").WithField("user_id", userID).WithField("dungeon_id", dungeonID)

	dungeon, err := model.FindDungeon(c, svr.db, dungeonID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.GinHandleError(c, log, http.StatusNotFound, err, "dungeon not found")
		} else {
			utils.GinHandleError(c, log, http.StatusInternalServerError, err, "failed to find dungeon")
		}
		return
	}
	if dungeon == nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "got nil dungeon")
		return
	}

	books, err := dungeon.GetBookIDs(c, svr.db)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "failed to fetch dungeon books")
		return
	}

	c.JSON(http.StatusOK, books)
}

// GetDungeonTagsDetail handles fetching the tags of a specific dungeon
// @Summary Get the tags of a specific dungeon
// @Description 获取复习计划的 TagNames
// @Tags dungeon
// @Produce json
// @Param id path uint64 true "Dungeon ID"
// @Success 200 {array} utils.UInt64
// @Failure 404 {object} utils.ErrorResponse "Dungeon not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /dungeon/dungeons/{id}/tags [get]
func (svr *Service) GetDungeonTagsDetail(c *gin.Context) {
	userID, dungeonID := utils.GinMustGetUserID(c), utils.GinMustGetID(c)
	log := wlog.ByCtx(c, "GetDungeonTagsDetail").WithField("user_id", userID).WithField("dungeon_id", dungeonID)

	dungeon, err := model.FindDungeon(c, svr.db, dungeonID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.GinHandleError(c, log, http.StatusNotFound, err, "dungeon not found")
		} else {
			utils.GinHandleError(c, log, http.StatusInternalServerError, err, "failed to find dungeon")
		}
		return
	}
	if dungeon == nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "got nil dungeon")
		return
	}

	tags, err := dungeon.GetTagIDs(c, svr.db)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to fetch dungeon tags")
		return
	}

	c.JSON(http.StatusOK, tags)
}

// GetDungeonItemsDetail handles fetching the items of a specific dungeon
// @Summary Get the items of a specific dungeon
// @Description 获取复习计划的 Items
// @Tags dungeon
// @Produce json
// @Param id path uint64 true "Dungeon ID"
// @Success 200 {array} utils.UInt64
// @Failure 404 {object} utils.ErrorResponse "Dungeon not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /dungeon/dungeons/{id}/items [get]
func (svr *Service) GetDungeonItemsDetail(c *gin.Context) {
	userID, dungeonID := utils.GinMustGetUserID(c), utils.GinMustGetID(c)
	log := wlog.ByCtx(c, "GetDungeonItemsDetail").WithField("user_id", userID).WithField("dungeon_id", dungeonID)

	dungeon, err := model.FindDungeon(c, svr.db, dungeonID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.GinHandleError(c, log, http.StatusNotFound, err, "dungeon not found")
		} else {
			utils.GinHandleError(c, log, http.StatusInternalServerError, err, "failed to find dungeon")
		}
		return
	}
	if dungeon == nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "got nil dungeon")
		return
	}

	pager := utils.GinGetPagerFromQuery(c)
	monsters, err := dungeon.GetDirectMonsters(svr.db, pager.Offset, pager.Limit)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to fetch dungeon monsters")
		return
	}

	resp := new(dto.RespMonsterList)
	for _, dm := range monsters {
		resp = resp.Append(new(dto.DungeonMonster).FromModel(dm))
	}
	resp.WithPager(pager).Response(c, "got dungeon monsters")
}
