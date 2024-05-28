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
// @Param id path uint64 true "Dungeon ID"
// @Param sort_by query string true "Sort by field (familiarity, difficulty, importance)"
// @Param offset query int true "Offset for pagination"
// @Param limit query int true "Limit for pagination"
// @Success 200 {object} dto.RespMonsterList
// @Failure 404 {object} utils.ErrorResponse "Dungeon not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /dungeon/campaigns/{id}/monsters [get]
func (svr *Service) GetMonstersOfCampaignDungeon(c *gin.Context) {
	id := c.Param("id")
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
	if err = svr.db.Where("id = ?", id).First(&dungeon).Error; err != nil {
		utils.GinHandleError(c, wlog.ByCtx(c), http.StatusNotFound, err, "Dungeon not found")
		return
	}

	monsters, err := dungeon.GetMonsters(svr.db, sortBy, offset, limit)
	if err != nil {
		utils.GinHandleError(c, wlog.ByCtx(c), http.StatusInternalServerError, err, "Failed to fetch dungeon monsters")
		return
	}

	resp := dto.RespMonsterList{
		Message: "monsters found",
		Data:    make([]*dto.Monster, 0, len(monsters)),
	}

	for _, monster := range monsters {
		resp.Append((&dto.Monster{}).FromModel(monster))
	}

	c.JSON(http.StatusOK, resp)
}

// GetNextMonstersOfCampaignDungeon handles fetching the next n monsters of a specific dungeon
// @Summary Get the next n monsters of a specific dungeon
// @Description 获取复习计划的后n个Monsters
// @TagNames dungeon
// @Produce json
// @Param id path uint64 true "Dungeon ID"
// @Param count query int true "Number of monsters to fetch"
// @Param sort_by query string true "Sort by field (familiarity, difficulty, importance)"
// @Success 200 {array} dto.RespMonsterList
// @Failure 404 {object} utils.ErrorResponse "Dungeon not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /dungeon/campaigns/{id}/next_monsters [get]
func (svr *Service) GetNextMonstersOfCampaignDungeon(c *gin.Context) {
	// todo: 定义这个规则，现在只是 demo

	id := c.Param("id")
	countStr := c.Query("count")
	sortBy := c.Query("sort_by")

	count, err := strconv.Atoi(countStr)
	if err != nil {
		count = 0
	}

	var dungeon model.Dungeon
	if err = svr.db.Where("id = ?", id).First(&dungeon).Error; err != nil {
		utils.GinHandleError(c, wlog.ByCtx(c), http.StatusNotFound, err, "Dungeon not found")
		return
	}

	monsters, err := dungeon.GetMonsters(svr.db, sortBy, 0, count)
	if err != nil {
		utils.GinHandleError(c, wlog.ByCtx(c), http.StatusInternalServerError, err, "Failed to fetch dungeon monsters")
		return
	}

	resp := dto.RespMonsterList{
		Message: "monsters found",
		Data:    make([]*dto.Monster, 0, len(monsters)),
	}

	for _, monster := range monsters {
		resp.Append((&dto.Monster{}).FromModel(monster))
	}

	c.JSON(http.StatusOK, resp)
}

// ReportCampaignResult handles reporting the result of a specific monster
// @Summary Report the result of a specific monster
// @Description 上报复习计划的Monster结果
// @TagNames dungeon
// @Accept json
// @Produce json
// @Param id path uint64 true "Dungeon ID"
// @Param result body ReqReportMonsterResult true "UserMonster result data"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse "Invalid request body"
// @Failure 404 {object} utils.ErrorResponse "Dungeon or UserMonster not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /dungeon/campaigns/{id}/report_result [post]
func (svr *Service) ReportCampaignResult(c *gin.Context) {
	var req ReqReportMonsterResult
	id := c.Param("id")

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinHandleError(c, wlog.ByCtx(c), http.StatusBadRequest, err, "Invalid request body")
		return
	}

	// 查询Dungeon记录
	var dungeon model.Dungeon
	if err := svr.db.Where("id = ?", id).First(&dungeon).Error; err != nil {
		utils.GinHandleError(c, wlog.ByCtx(c), http.StatusNotFound, err, "Dungeon not found")
		return
	}

	// 查询Monster记录
	var item model.Item
	if err := svr.db.Where("id = ?", req.MonsterID).First(&item).Error; err != nil {
		utils.GinHandleError(c, wlog.ByCtx(c), http.StatusNotFound, err, "UserMonster not found")
		return
	}

	// todo: 处理Monster结果 - 根据需求调整处理逻辑，例如更新Monster的熟练度或状态等

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "UserMonster result reported",
	})
}

// GetCampaignDungeonTodayConclusion handles fetching the results of a specific dungeon
// @Summary Get the results of a specific dungeon
// @Description 获取复习计划的结果
// @TagNames dungeon
// @Produce json
// @Param id path uint64 true "Dungeon ID"
// @Success 200 {object} RespDungeonResults
// @Failure 404 {object} utils.ErrorResponse "Dungeon not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /dungeon/campaigns/{id}/results [get]
func (svr *Service) GetCampaignDungeonTodayConclusion(c *gin.Context) {
	// log := wlog.ByCtx(c, "GetCampaignDungeonTodayConclusion")
	var dungeon model.Dungeon
	id := c.Param("id")

	if err := svr.db.Where("id = ?", id).First(&dungeon).Error; err != nil {
		utils.GinHandleError(c, wlog.ByCtx(c), http.StatusNotFound, err, "Dungeon not found")
		return
	}
	//
	//// 获取所有Monster
	//monsters, err := model.GetDungeonMonsters(svr.db, dungeon.ID, "", 0, 9999)
	//if err != nil {
	//	utils.GinHandleError(c, wlog.ByCtx(c), http.StatusInternalServerError, err, "Failed to fetch dungeon items")
	//	return
	//}
	//
	//// todo: 计算今天的复习难度 - 在实际应用中，这里可以根据需求定义复杂的计算规则
	//todayDifficulty := 0
	//for _, monster := range monsters {
	//	// demo: 根据Monster的状态或其他信息来计算难度, 这里假设每个Monster的难度为1
	//	log.Infof("got monster %v", monster)
	//	todayDifficulty += 1
	////}
	//
	//resp := RespDungeonResults{
	//	TotalMonsters:   len(monsters),
	//	TodayDifficulty: todayDifficulty,
	//}

	// todo

	c.JSON(http.StatusOK, "")
}
