package dao

import (
	"context"
	"time"

	"github.com/appootb/protobuf/go/service"
	"github.com/appootb/substratum/metadata"
	"github.com/appootb/substratum/storage"
	"github.com/jinzhu/gorm"
)

type Option func(*Base)

func WithContext(ctx context.Context) Option {
	return func(base *Base) {
		component := service.ComponentNameFromContext(ctx)
		db := storage.ContextStorage(ctx, component).GetDB()
		if md := metadata.RequestMetadata(ctx); md != nil && md.GetIsDebug() {
			db = db.Debug()
		}
		base.tx = db
	}
}

func WithDB(tx *gorm.DB) Option {
	return func(base *Base) {
		base.tx = tx
	}
}

func WithPrimaryKey(id uint64) Option {
	return func(base *Base) {
		base.ID = id
	}
}

func WithCreatedTime(timestamp time.Time) Option {
	return func(base *Base) {
		base.CreatedAt = timestamp
	}
}
