package model

import (
	"context"
	"errors"
	"strings"

	"github.com/bagaking/goulp/wlog"
	"github.com/bagaking/memorianexus/internal/util"
	"github.com/khicago/irr"
	"gorm.io/gorm"
)

type Tag struct {
	ID   util.UInt64 `gorm:"primaryKey;autoIncrement:false"`
	Name string      `gorm:"unique;not null"`
}

var ErrInvalidTagName = errors.New("invalid tag name")

func FindOrUpdateTagByName(c context.Context, tx *gorm.DB, tagName string, id util.UInt64) (*Tag, error) {
	log := wlog.ByCtx(c, "FindOrUpdateTagByName")
	tagName = strings.TrimSpace(tagName)
	if tagName == "" {
		return nil, ErrInvalidTagName
	}
	// 查找或创建 Tag
	tag := &Tag{}
	err := tx.Where(&Tag{Name: tagName}).First(tag).Error
	if err == nil {
		return tag, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, irr.Wrap(err, "find tag failed")
	}
	tag = &Tag{ID: id, Name: tagName}
	if err = tx.Create(tag).Error; err != nil {
		return nil, irr.Wrap(err, "create tag failed")
	}
	log.Infof("new tag %v create, id_set= %d", tag, id)
	return tag, nil
}
