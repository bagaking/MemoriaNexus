package item

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bagaking/ankibuild/anki"
	"github.com/khicago/irr"

	"github.com/bagaking/goulp/wlog"
	"github.com/gin-gonic/gin"
	"github.com/khicago/got/util/typer"

	"github.com/bagaking/memorianexus/internal/utils"
	"github.com/bagaking/memorianexus/src/def"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/bagaking/memorianexus/src/module/dto"
)

// UploadItems handles uploading a file to create multiple items.
// @Summary Upload items from a file
// @Description Upload a file to create multiple items in the system.
// @Tags item
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File containing items data, support csv and toml file"
// @Param book_id query string false "Book ID"
// @Success 201 {object} dto.RespItemList "Successfully created items from file"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Router /items/upload [post]
func (svr *Service) UploadItems(c *gin.Context) {
	userID := utils.GinMustGetUserID(c)
	log := wlog.ByCtx(c, "UploadItems").WithField("user_id", userID)

	var req ReqUploadItems
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, err, "invalid query parameters")
		return
	}

	var book *model.Book
	if req.BookID != nil {
		b, err := model.FindBook(c, svr.db, *req.BookID)
		if err != nil {
			utils.GinHandleError(c, log, http.StatusInternalServerError, err, "failed to find the book")
			return
		}
		book = b
	}

	file, err := c.FormFile("file")
	if err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, err, "failed to get file")
		return
	}

	// 打开上传的文件
	f, err := file.Open()
	if err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "failed to open file")
		return
	}
	defer f.Close()

	// 解析文件内容
	items, itemTagRef, err := parseItemsFromFile(c, f, file.Filename)
	if err != nil {
		utils.GinHandleError(c, log, http.StatusBadRequest, err, "failed to parse file")
		return
	}
	ids, err := utils.MGenIDU64(c, len(items))
	if err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "failed to generate item id")
		return
	}
	for i := range items {
		item := items[i]
		item.CreatorID = userID
		item.ID = ids[i]
		item.CreatedAt = time.Now()
	}

	// 保存解析后的学习材料
	if err = model.CreateItems(c, svr.db, userID, items, itemTagRef); err != nil {
		utils.GinHandleError(c, log, http.StatusInternalServerError, err, "failed to save items")
		return
	}

	if book != nil {
		successItemIDs, err := book.MPutItems(c, svr.db, typer.SliceMap(items, func(from *model.Item) utils.UInt64 {
			return from.ID
		}))
		if err != nil {
			log.WithError(err).Errorf("mput items failed, success= %v", successItemIDs)
			utils.GinHandleError(c, log, http.StatusInternalServerError, err, "failed to save items")
			return
		}
		log.Infof("mput items successfully, success= %v", successItemIDs)
	}

	new(dto.RespItemList).Append(typer.SliceMap(items, func(from *model.Item) *dto.Item {
		return new(dto.Item).FromModel(from)
	})...).Response(c, "items created from file")
}

func parseItemsFromFile(ctx context.Context, r io.Reader, filename string) ([]*model.Item, map[*model.Item][]string, error) {
	ext := filepath.Ext(filename)
	switch ext {
	case ".csv":
		return parseItemsFromCSV(ctx, r)
	case ".toml":
		content, err := io.ReadAll(r)
		if err != nil {
			return nil, nil, err
		}
		return parseItemsFromTOML(ctx, content)
	default:
		return nil, nil, errors.New("unsupported file format")
	}
}

func parseItemsFromCSV(ctx context.Context, r io.Reader) ([]*model.Item, map[*model.Item][]string, error) {
	reader := csv.NewReader(r)
	var items []*model.Item
	itemTagRef := make(map[*model.Item][]string)

	// 读取 CSV 文件的每一行
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}

		// 假设 CSV 文件的每一行包含以下字段：Type, Content, Difficulty, Importance, Tags
		if len(record) < 5 {
			return nil, nil, irr.Error("invalid record length")
		}

		difficulty, err := strconv.ParseUint(record[2], 16, 8)
		if err != nil {
			return nil, nil, err
		}

		importance, err := strconv.ParseUint(record[3], 16, 8)
		if err != nil {
			return nil, nil, err
		}

		tags := strings.Split(record[4], ",")

		item := &model.Item{
			Type:       record[0],
			Content:    record[1],
			Difficulty: def.DifficultyLevel(difficulty),
			Importance: def.ImportanceLevel(importance),
		}
		items = append(items, item)
		itemTagRef[item] = tags
	}

	return items, itemTagRef, nil
}

func parseItemsFromTOML(ctx context.Context, content []byte) ([]*model.Item, map[*model.Item][]string, error) {
	barn, err := anki.ParseTomlContent(ctx, content)
	if err != nil {
		return nil, nil, err
	}

	itemTagRef := make(map[*model.Item][]string)

	var items []*model.Item

	tags := barn.Tags
	for _, card := range barn.QnAs {
		item := &model.Item{
			Type:    model.TyItemFlashCard,
			Content: "## " + card.Question + "\n\n" + card.Answer,
		}
		items = append(items, item)
		itemTagRef[item] = append(tags, card.Tags...)
	}

	return items, itemTagRef, nil
}
