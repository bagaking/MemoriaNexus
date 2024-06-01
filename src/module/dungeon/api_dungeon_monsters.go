package dungeon

import (
	"net/http"
	"strconv"

	"github.com/bagaking/memorianexus/src/module/dto"

	"github.com/bagaking/goulp/wlog"
	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/gin-gonic/gin"
)

type ReqReportMonsterResult struct {
	MonsterID uint64 `json:"monster_id"`
	Result    string `json:"result"` // "unknown", "familiar", "remembered"
}

// GetMonstersOfCampaignDungeon handles fetching all the monsters of a specific campaign dungeon
// @Summary Get all the monsters of a specific campaign dungeon
// @Description 获取复习计划的所有Monsters
// @TagNames dungeon
// @Produce json
// @Param id path string true "Dungeon ID"
// @Param sort_by query string true "Sort by field (familiarity, difficulty, importance)"
// @Param offset query int true "Offset for pagination"
// @Param limit query int true "Limit for pagination"
// @Success 200 {object} dto.RespMonsterList
// @Failure 404 {object} utils.ErrorResponse "Dungeon not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /dungeon/campaigns/{id}/monsters [get]
func (svr *Service) GetMonstersOfCampaignDungeon(c *gin.Context) {
	log := wlog.ByCtx(c, "GetMonstersOfCampaignDungeon")
	dungeonID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, err, "Invalid dungeon ID")
		return
	}

	sortBy := c.DefaultQuery("sort_by", "familiarity")
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "10")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	var dungeon model.Dungeon
	if err = svr.db.Where("id = ?", dungeonID).First(&dungeon).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusNotFound, err, "Dungeon not found")
		return
	}

	monsters, err := dungeon.GetMonsters(svr.db, sortBy, offset, limit)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to fetch dungeon monsters")
		return
	}

	resp := new(dto.RespMonsterList)
	for _, monster := range monsters {
		resp.Append(new(dto.DungeonMonster).FromModel(monster))
	}

	resp.Response(c)
}

// GetNextMonstersOfCampaignDungeon handles fetching the next n monsters of a specific dungeon
// @Summary Get the next n monsters of a specific dungeon
// @Description 获取复习计划的后n个Monsters
// @TagNames dungeon
// @Produce json
// @Param id path string true "Dungeon ID"
// @Param count query int true "Number of monsters to fetch"
// @Param sort_by query string true "Sort by field (familiarity, difficulty, importance)"
// @Success 200 {object} dto.RespMonsterList
// @Failure 404 {object} utils.ErrorResponse "Dungeon not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /dungeon/campaigns/{id}/next_monsters [get]
func (svr *Service) GetNextMonstersOfCampaignDungeon(c *gin.Context) {
	log := wlog.ByCtx(c, "GetNextMonstersOfCampaignDungeon")
	dungeonID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, err, "Invalid dungeon ID")
		return
	}

	countStr := c.Query("count")
	sortBy := c.Query("sort_by")

	count, err := strconv.Atoi(countStr)
	if err != nil {
		count = 0
	}

	var dungeon model.Dungeon
	if err = svr.db.Where("id = ?", dungeonID).First(&dungeon).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusNotFound, err, "Dungeon not found")
		return
	}

	monsters, err := dungeon.GetMonsters(svr.db, sortBy, 0, count)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to fetch dungeon monsters")
		return
	}

	resp := new(dto.RespMonsterList)
	for _, monster := range monsters {
		resp.Append(new(dto.DungeonMonster).FromModel(monster))
	}
	resp.Response(c)
}
