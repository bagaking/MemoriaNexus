package book

import (
	"context"
	"errors"
	"github.com/bagaking/goulp/wlog"
	"github.com/khicago/irr"
	"gorm.io/gorm"

	"github.com/bagaking/memorianexus/internal/util"
	"github.com/bagaking/memorianexus/src/model"
)

// updateBookTagsRef updates the tags association for a book in the database.
func updateBookTagsRef(ctx context.Context, tx *gorm.DB, bookID util.UInt64, tags []string) error {
	log := wlog.ByCtx(ctx, "updateBookTagsRef")
	tagIDs, err := util.MGenIDU64(ctx, len(tags))
	if err != nil {
		return irr.Wrap(err, "failed to generate IDs for tags")
	}

	// Delete existing associations for this book.
	if err := tx.Where("book_id = ?", bookID).Delete(&model.BookTag{}).Error; err != nil {
		log.WithError(err).Error("Failed to delete existing book tags")
		return irr.Wrap(err, "failed to delete existing book tags")
	}

	for i, tagName := range tags {
		// Find or create the tag.
		tag, err := model.FindOrUpdateTagByName(ctx, tx, tagName, tagIDs[i])
		if err != nil {
			if errors.Is(err, model.ErrInvalidTagName) {
				continue // skip invalid tag names
			}
			log.WithError(err).Errorf("Failed to find or create tag '%s'", tagName)
			return irr.Wrap(err, "failed to find or create tag")
		}
		// Associate the tag with the book.
		bookTag := &model.BookTag{
			BookID: bookID,
			TagID:  tag.ID,
		}
		if err := tx.Create(bookTag).Error; err != nil {
			log.WithError(err).Error("Failed to associate tag with book")
			return irr.Wrap(err, "failed to associate tag with book")
		}
	}
	return nil
}
