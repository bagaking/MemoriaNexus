package book

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/bagaking/goulp/wlog"
	"github.com/gin-gonic/gin"
	"github.com/khicago/got/util/typer"
	"github.com/khicago/irr"

	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/bagaking/memorianexus/src/module/dto"
)

// GetItemsOfBook handles retrieving a list of books with pagination.
// @Summary Get item list of books with pagination
// @Description Get a paginated list of items for the book.
// @Tags book
// @Accept json
// @Produce json
// @Param page query int false "Page number for pagination" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} dto.RespItemList "items of the book found"
// @Router /books/{id}/items [get]
func (svr *Service) GetItemsOfBook(c *gin.Context) {
	userID := utils.GinMustGetUserID(c)
	bookID := utils.GinMustGetID(c)
	pager := utils.GinGetPagerFromQuery(c)
	log := wlog.ByCtx(c, "GetItemsOfBook").WithField("user_id", userID).WithField("pager", pager)

	items, err := model.GetItemsOfBook(svr.db, bookID, pager.Offset, pager.Limit)
	if err != nil {
		log.WithError(err).Errorf("Failed to fetch books for user %v", userID)
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "error when fetching book items")
	}

	resp := new(dto.RespItemList)
	for _, item := range items {
		resp.Append(new(dto.Item).FromModel(item))
	}
	log.Infof("get book items, offset= %v, limit= %v, items_len= %v", pager.Offset, pager.Limit, len(items))

	if err = svr.db.Model(&model.BookItem{}).Where(model.BookItem{BookID: bookID}).Count(&pager.Total).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warnf("try to count book items failed, %v", err)
		}
	}

	resp.Response(c, "items of the book found")
}

// AddItemsToBook handles adding items to a book
// @Summary Add items to a book
// @Description Add a list of items to the specified book.
// @Tags book
// @Accept json
// @Produce json
// @Param id path uint64 true "Book ID"
// @Param body body ReqAddItems true "List of item IDs to add"
// @Success 200 {object} dto.SuccessResponse "items added to book successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} utils.ErrorResponse "Failed to upsert book items"
// @Router /books/{id}/items [post]
func (svr *Service) AddItemsToBook(c *gin.Context) {
	userID := utils.GinMustGetUserID(c)
	bookID := utils.GinMustGetID(c)
	log := wlog.ByCtx(c, "AddItemsToBook").WithField("user_id", userID)

	var req ReqAddItems
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, err, "invalid request body")
		return
	}

	bookItems := typer.SliceMap(req.ItemIDs, func(from utils.UInt64) model.BookItem {
		return model.BookItem{
			BookID: bookID,
			ItemID: from,
		}
	})

	// 使用 OnConflict 方法处理 upsert 逻辑
	if err := svr.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "book_id"}, {Name: "item_id"}},
		DoNothing: true, // 如果冲突则不做任何操作
	}).Create(&bookItems).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "failed to upsert book items")
		return
	}

	new(dto.SuccessResponse).With(bookItems).Response(c, "items added to book successfully")
}

// RemoveItemsFromBook handles removing items from a book
// @Summary Remove items from a book
// @Description Remove a list of items from the specified book.
// @Tags book
// @Accept json
// @Produce json
// @Param id path uint64 true "Book ID"
// @Param item_ids query string true "Comma-separated list of item IDs to remove"
// @Success 200 {object} dto.SuccessResponse "items removed from book successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request parameters"
// @Failure 500 {object} utils.ErrorResponse "Failed to remove book items"
// @Router /books/{id}/items [delete]
func (svr *Service) RemoveItemsFromBook(c *gin.Context) {
	userID := utils.GinMustGetUserID(c)
	bookID := utils.GinMustGetID(c)
	log := wlog.ByCtx(c, "RemoveItemsFromBook").WithField("user_id", userID)

	itemIDsStr := c.Query("item_ids")
	if itemIDsStr == "" {
		utils.GinHandleError(c, log, http.StatusBadRequest, irr.Error("item_ids query parameter is required"), "invalid params")
		return
	}

	itemIDs := strings.Split(itemIDsStr, ",")
	var itemIDsUInt64 []utils.UInt64
	for _, idStr := range itemIDs {
		id := utils.UInt64(0)
		err := (&id).Scan(idStr)
		if err != nil {
			utils.GinHandleError(c, log, http.StatusBadRequest, err, fmt.Sprintf("invalid item_id: "+idStr))
			return
		}
		itemIDsUInt64 = append(itemIDsUInt64, id)
	}

	if err := svr.db.Where("book_id = ? AND item_id IN ( ? )", bookID, itemIDsUInt64).Delete(&model.BookItem{}).Error; err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, err, "failed to remove book items")
		return
	}

	new(dto.SuccessResponse).With(itemIDsUInt64).Response(c, "items removed from book successfully")
}
