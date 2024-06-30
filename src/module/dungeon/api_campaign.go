package dungeon

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/bagaking/memorianexus/src/def"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/khicago/got/util/typer"
	"github.com/khicago/irr"

	"github.com/bagaking/memorianexus/src/module/dto"

	"github.com/bagaking/goulp/wlog"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/gin-gonic/gin"
)

type ReqReportMonsterResult struct {
	MonsterID utils.UInt64     `json:"monster_id"`
	Result    def.AttackResult `json:"result"` // "defeat", "miss", "hit", "kill", "complete"
}

// GetMonstersOfCampaignDungeon handles fetching all the monsters of a specific campaign dungeon
// @Summary Get all the monsters of a specific campaign dungeon
// @Description 获取复习计划的所有Monsters
// @Tags dungeon
// @Produce json
// @Param id path uint64 true "Dungeon ID"
// @Param sort_by query string true "Sort by field (familiarity, difficulty, importance)"
// @Param page query int true "page for pagination"
// @Param limit query int true "Limit for pagination"
// @Success 200 {object} dto.RespMonsterList "Successfully retrieved monsters"
// @Failure 404 {object} utils.ErrorResponse "Dungeon not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /dungeon/campaigns/{id}/monsters [get]
func (svr *Service) GetMonstersOfCampaignDungeon(c *gin.Context) {
	userID, campaignID := utils.GinMustGetUserID(c), utils.GinMustGetID(c)
	log := wlog.ByCtx(c, "GetMonstersOfCampaignDungeon").WithField("user_id", userID).WithField("campaign_id", campaignID)

	offsetStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(offsetStr)
	if err != nil {
		page = 0
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	var dungeon model.Dungeon
	if err = svr.db.Where("id = ?", campaignID).First(&dungeon).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusNotFound, err, "Dungeon not found")
		return
	}

	pager := new(utils.Pager).SetPageAndLimit(page, limit)
	monsters, err := dungeon.GetMonsters(svr.db, pager.Offset, pager.Limit)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to fetch dungeon monsters")
		return
	}

	resp := new(dto.RespMonsterList).WithPager(pager)
	for _, monster := range monsters {
		resp.Append(new(dto.DungeonMonster).FromModel(monster))
	}
	resp.Response(c)
}

// GetMonstersForCampaignPractice handles fetching monsters for review
// @Summary Get monsters for review
// @Description 从 Campaign Dungeon 中提取一些要复习的 Monster 缓存到本地
// @Tags dungeon
// @Produce json
// @Param id path uint64 true "Dungeon ID"
// @Param count query int true "Number of monsters to fetch"
// @Success 200 {object} dto.RespMonsterList "Successfully retrieved monsters"
// @Failure 404 {object} utils.ErrorResponse "Dungeon not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /dungeon/campaigns/{id}/practice [get]
func (svr *Service) GetMonstersForCampaignPractice(c *gin.Context) {
	userID, campaignID := utils.GinMustGetUserID(c), utils.GinMustGetID(c)
	l, ctx := wlog.ByCtxAndCache(c, "GetMonstersForCampaignPractice")
	log := l.WithField("user_id", userID).WithField("campaign_id", campaignID)

	countStr := c.DefaultQuery("count", "20")
	strategy := c.DefaultQuery("strategy", "classic")

	count, err := strconv.Atoi(countStr)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, irr.Wrap(err, "count= %v", countStr), "invalid argument")
	}
	log = log.WithField("count", count)
	pager := new(utils.Pager).SetFirstCount(count)

	var dungeon model.Dungeon
	if err = svr.db.Where("id = ?", campaignID).First(&dungeon).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusNotFound, err, "dungeon not found")
		return
	}

	monsters, err := dungeon.GetMonstersForPractice(ctx, svr.db, strategy, pager.Limit)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "failed to fetch dungeon monsters")
		return
	}

	resp := new(dto.RespMonsterList).WithPager(pager)
	dtoMonsters := typer.SliceMap(monsters,
		func(from model.DungeonMonster) *dto.DungeonMonster {
			return new(dto.DungeonMonster).FromModel(from)
		})
	log.Infof("got monsters= %v", dtoMonsters)
	resp.Append(dtoMonsters...).Response(c)
}

// SubmitCampaignResult handles reporting the result of a specific monster recall
// @Summary Report the result of a specific monster recall
// @Description 上报复习计划的Monster结果
// @Tags dungeon
// @Accept json
// @Produce json
// @Param id path uint64 true "Dungeon ID"
// @Param result body ReqReportMonsterResult true "UserMonster result data"
// @Success 200 {object} dto.SuccessResponse "Successfully reported result"
// @Failure 400 {object} utils.ErrorResponse "Invalid request body"
// @Failure 404 {object} utils.ErrorResponse "Dungeon or UserMonster not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /dungeon/campaigns/{id}/submit [post]
func (svr *Service) SubmitCampaignResult(c *gin.Context) {
	userID, campaignID := utils.GinMustGetUserID(c), utils.GinMustGetID(c)
	log := wlog.ByCtx(c, "SubmitCampaignResult").WithField("user_id", userID).WithField("campaign_id", campaignID)

	var req ReqReportMonsterResult
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, err, "invalid request body")
		return
	}

	// 更新DungeonMonster的显影程度和下次复习时间
	var dm model.DungeonMonster
	if err := svr.db.Where("dungeon_id = ? AND item_id = ?", campaignID, req.MonsterID).First(&dm).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusNotFound, err, "DungeonMonster not found")
		return
	}
	log = log.WithField("item_id", dm.ItemID)

	// 处理Monster结果 - 根据需求调整处理逻辑，例如更新Monster的熟练度或状态等
	damageRate := req.Result.DamageRate()
	if damageRate <= 0 {
		utils.GinHandleError(c, log, http.StatusBadRequest, irr.Error("invalid attack result %s", req.Result), "Invalid result")
		return
	}

	// 更新UserMonster的熟练度
	newFamiliarity := CalculateNewFamiliarity(dm.Familiarity, damageRate, dm.PracticeAt, dm.Difficulty)
	userMonster := model.UserMonster{
		UserID:      userID,
		ItemID:      req.MonsterID,
		Familiarity: newFamiliarity,
	}
	if err := svr.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "item_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"familiarity"}),
	}).Create(&userMonster).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "failed to update UserMonster familiarity")
		return
	}
	log.Infof("damage calculate, last_practice_at %v, damage_rate= %v, difficulty= %v, current= %v, new= %v", dm.PracticeAt, damageRate, dm.Difficulty, dm.Familiarity, newFamiliarity)

	// 计算下次复习时间
	var userSettings model.ProfileMemorizationSetting
	if err := svr.db.Where("id = ?", userID).First(&userSettings).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.GinHandleError(c, log, http.StatusInternalServerError, err, "failed to fetch user settings")
			return
		}
		userSettings = model.DefaultMemorizationSetting
		userSettings.ID = userID
	}

	nextRecallTime := CalculateNextPracticeAt(c, newFamiliarity, dm.Importance, &userSettings)
	updater := map[string]any{
		"visibility":       utils.Percentage(newFamiliarity.Times(dm.Visibility.NormalizedFloat())),
		"familiarity":      newFamiliarity,
		"practice_at":      time.Now(),
		"next_practice_at": nextRecallTime,
		"practice_count":   gorm.Expr("practice_count + ?", 1),
	}

	if err := svr.db.Model(&dm).
		Where("dungeon_id = ? AND item_id = ?", campaignID, req.MonsterID).
		Updates(updater).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "failed to update DungeonMonster visibility and next recall time")
		return
	}
	log.Infof("next_practice_at updated, last_practice_at= %v, new_familiarity= %v, importance= %v, next_recall_at= %v", dm.PracticeAt, newFamiliarity, dm.Importance, nextRecallTime)

	new(dto.RespMonsterUpdate).With(
		dto.Updater[*dto.DungeonMonster]{
			From:    new(dto.DungeonMonster).FromModel(dm),
			Updates: updater,
		}).Response(c, "user-monster practice result updated")
}

func (svr *Service) GetCampaignDungeonConclusionOfToday(c *gin.Context) {
	userID, campaignID := utils.GinMustGetUserID(c), utils.GinMustGetID(c)
	log := wlog.ByCtx(c, "GetCampaignDungeonConclusionOfToday").WithField("user_id", userID).WithField("campaign_id", campaignID)

	var dungeon model.Dungeon
	if err := svr.db.Where("id = ?", campaignID).First(&dungeon).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusNotFound, err, "Dungeon not found")
		return
	}
	//
	//// 获取所有Monster
	//monsters, err := model.GetDirectMonsters(svr.db, dungeon.ID, "", 0, 9999)
	//if err != nil {
	//	utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to fetch dungeon items")
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
