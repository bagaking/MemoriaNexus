package dungeon

import (
	"net/http"

	"github.com/bagaking/goulp/wlog"
	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/bagaking/memorianexus/src/module/dto"
	"github.com/gin-gonic/gin"
)

type RespCampaignResult struct {
	TotalMonsters   int `json:"total_monsters"`
	TodayDifficulty int `json:"today_difficulty"`
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
// @Success 200 {object} RespCampaignResult
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
	//resp := RespCampaignResult{
	//	TotalMonsters:   len(monsters),
	//	TodayDifficulty: todayDifficulty,
	//}
	// todo
	c.JSON(http.StatusOK, "")
}
