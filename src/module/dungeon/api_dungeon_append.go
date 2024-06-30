package dungeon

import (
	"errors"
	"net/http"

	"gorm.io/gorm"

	"github.com/bagaking/goulp/wlog"
	"github.com/gin-gonic/gin"
	"github.com/khicago/irr"

	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/bagaking/memorianexus/src/module/dto"
)

// ReqAddDungeonBooks defines the request structure for adding books to a dungeon
type ReqAddDungeonBooks struct {
	Books []utils.UInt64 `json:"books"`
}

// ReqAddDungeonItems defines the request structure for adding items to a dungeon
type ReqAddDungeonItems struct {
	Items []utils.UInt64 `json:"items"`
}

// ReqAddDungeonTags defines the request structure for adding tags to a dungeon
type ReqAddDungeonTags struct {
	TagNames []string `json:"tag_names"`
}

// AppendBooksToDungeon handles adding books to an existing dungeon
// @Summary Add books to an existing dungeon
// @Description 向现有复习计划添加书籍
// @Tags dungeon
// @Accept json
// @Produce json
// @Param id path string true "Dungeon ID"
// @Param books body ReqAddDungeonBooks true "Books to add"
// @Success 200 {object} dto.RespDungeon
// @Failure 400 {object} utils.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /dungeon/dungeons/{id}/books [post]
func (svr *Service) AppendBooksToDungeon(c *gin.Context) {
	userID, dungeonID := utils.GinMustGetUserID(c), utils.GinMustGetID(c)
	log := wlog.ByCtx(c, "AppendBooksToDungeon").WithField("user_id", userID).WithField("dungeon_id", dungeonID)

	var req ReqAddDungeonBooks
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, irr.Wrap(err, "parse request body failed"), "Invalid request body")
		return
	}

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
	if dungeon.UserID != userID {
		utils.GinHandleError(c, log, http.StatusForbidden, err, "got permission denied")
	}

	if err = dungeon.AddMonster(c, svr.db, model.MonsterSourceBook, req.Books); err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to add books to dungeon")
		return
	}

	new(dto.RespDungeon).With(new(dto.Dungeon).FromModel(dungeon)).Response(c)
}

// AppendItemsToDungeon handles adding items to an existing dungeon
// @Summary Add items to an existing dungeon
// @Description 向现有复习计划添加学习材料
// @Tags dungeon
// @Accept json
// @Produce json
// @Param id path string true "Dungeon ID"
// @Param items body ReqAddDungeonItems true "Items to add"
// @Success 200 {object} dto.RespDungeon
// @Failure 400 {object} utils.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /dungeon/dungeons/{id}/items [post]
func (svr *Service) AppendItemsToDungeon(c *gin.Context) {
	userID, dungeonID := utils.GinMustGetUserID(c), utils.GinMustGetID(c)
	log := wlog.ByCtx(c, "AppendItemsToDungeon").WithField("user_id", userID).WithField("dungeon_id", dungeonID)

	var req ReqAddDungeonItems
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, irr.Wrap(err, "parse request body failed"), "Invalid request body")
		return
	}

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
	if dungeon.UserID != userID {
		utils.GinHandleError(c, log, http.StatusForbidden, err, "got permission denied")
	}

	if err = dungeon.AddMonster(c, svr.db, model.MonsterSourceItem, req.Items); err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "failed to add items to dungeon")
		return
	}

	new(dto.RespDungeon).With(new(dto.Dungeon).FromModel(dungeon)).Response(c)
}

// AppendTagsToDungeon handles adding tags to an existing dungeon
// @Summary Add tags to an existing dungeon
// @Description 向现有复习计划添加标签
// @Tags dungeon
// @Accept json
// @Produce json
// @Param id path string true "Dungeon ID"
// @Param tags body ReqAddDungeonTags true "Tags to add"
// @Success 200 {object} dto.RespDungeon
// @Failure 400 {object} utils.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /dungeon/dungeons/{id}/tags [post]
func (svr *Service) AppendTagsToDungeon(c *gin.Context) {
	userID, dungeonID := utils.GinMustGetUserID(c), utils.GinMustGetID(c)
	log := wlog.ByCtx(c, "AppendTagsToDungeon").WithField("user_id", userID).WithField("dungeon_id", dungeonID)

	var req ReqAddDungeonTags
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, irr.Wrap(err, "parse request body failed"), "Invalid request body")
		return
	}

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
	if dungeon.UserID != userID {
		utils.GinHandleError(c, log, http.StatusForbidden, err, "got permission denied")
	}

	tagIDs, err := model.FindTagsIDByName(svr.db, req.TagNames)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to fetch tag IDs")
		return
	}

	if err = dungeon.AddMonster(c, svr.db, model.MonsterSourceTag, tagIDs); err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to add tags to dungeon")
		return
	}

	new(dto.RespDungeon).With(new(dto.Dungeon).FromModel(dungeon)).Response(c)
}
