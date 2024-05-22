package item

import (
	"context"
	"errors"
	"github.com/bagaking/goulp/wlog"
	"github.com/bagaking/memorianexus/internal/util"
	"github.com/bagaking/memorianexus/src/model"
	"github.com/khicago/irr"
	"gorm.io/gorm"
)

// updateItemTagsRef 处理 Tags 更新
func updateItemTagsRef(ctx context.Context, tx *gorm.DB, itemID util.UInt64, tags []string) error {
	log := wlog.ByCtx(ctx, "updateItemTagsRef")
	tagIDs, err := util.MGenIDU64(ctx, len(tags))
	log.Infof("tags %v ids generated: %v", tags, tagIDs)
	if err != nil {
		return irr.Wrap(err, "generate id for tags failed")
	}

	// todo: 删除现有的 tags, 先粗暴清除当前 Item 的所有 Tag 关联 (这个简易实现里面，删除涉及全表了，先打一版)
	dmlDelete := tx.Where("item_id = ?", itemID).Delete(model.ItemTag{})
	if err = dmlDelete.Error; err != nil {
		return irr.Wrap(err, "delete item tags failed")
	}
	log.Infof("tags %v ids dropped count %v", tags, dmlDelete.RowsAffected)

	// 为 Item 添加新的 Tag 关联
	for i, tagName := range tags {
		tag, errTag := model.FindOrUpdateTagByName(ctx, tx, tagName, tagIDs[i])
		if errTag != nil {
			if errors.Is(errTag, model.ErrInvalidTagName) {
				continue
			}
			return irr.Wrap(err, "upsert tag failed")
		}
		itemTag := &model.ItemTag{
			ItemID: itemID,
			TagID:  tag.ID,
		}
		if err = tx.FirstOrCreate(itemTag).Error; err != nil {
			return irr.Wrap(err, "upsert item_tag_ref failed")
		}
	}
	return nil
}
