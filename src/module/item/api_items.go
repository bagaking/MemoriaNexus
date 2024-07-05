package item

import (
	"net/http"

	"github.com/bagaking/goulp/wlog"

	"github.com/gin-gonic/gin"

	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/bagaking/memorianexus/src/module/dto"
)

// GetItems handles retrieving a list of items with optional filters and pagination.
// @Summary Get a list of items with optional filters
// @Description Get a list of items for the user with optional filters for book and type and support for pagination.
// @Tags item
// @Accept json
// @Produce json
// @Param user_id query uint64 false "User ID"
// @Param type query string false "Type of item"
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} dto.RespItemList "Successfully retrieved items"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Router /items [get]
func (svr *Service) GetItems(c *gin.Context) {
	userID := utils.GinMustGetUserID(c)
	pager := utils.GinGetPagerFromQuery(c)
	log := wlog.ByCtx(c, "GetItems").WithField("user_id", userID).WithField("pager", pager)

	var req ReqGetItems
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, err, "Invalid query parameters")
		return
	}

	if req.UserID <= 0 { // 如果不指定用户，搜索的就是自己的
		req.UserID = userID
	}
	query := &model.Item{
		CreatorID: req.UserID,
		Type:      req.Type,
	}
	var items []model.Item
	if err := svr.db.Model(query).Where(query).Offset(pager.Offset).Limit(pager.Limit).Find(&items).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to retrieve items")
		return
	}

	// todo: cache this
	var total int64
	if err := svr.db.Model(query).Where(query).Count(&total).Error; err != nil {
		log.WithError(err).Warnf("Failed to count items")
		return
	}
	pager = pager.SetTotal(total)

	resp := new(dto.RespItemList)
	// 转换 Item 为 Item
	for _, item := range items {
		tags, err := model.GetTagsByEntity(c, svr.db, item.ID)
		if err != nil {
			utils.GinHandleError(c, log, http.StatusInternalServerError, err, "Failed to get item tag names")
			return
		}
		resp.Append(new(dto.Item).FromModel(&item, tags...))
	}
	resp.WithPager(pager).Response(c, "items found")
}
