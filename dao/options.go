package dao

import (
	"context"
	"time"

	"github.com/appootb/substratum/metadata"
	"github.com/appootb/substratum/service"
	"github.com/appootb/substratum/storage"
	"gorm.io/gorm"
)

type Option func(*Base)

func WithContext(ctx context.Context) Option {
	return func(base *Base) {
		component := service.ComponentNameFromContext(ctx)
		base.rw = storage.ContextStorage(ctx, component).GetDB()
		base.ro = storage.ContextStorage(ctx, component).GetDB(true)
		if md := metadata.IncomingMetadata(ctx); md != nil && md.GetIsDebug() {
			base.rw = base.rw.Debug()
			base.ro = base.ro.Debug()
		}
	}
}

func WithDB(tx *gorm.DB) Option {
	return func(base *Base) {
		base.rw = tx
		base.ro = tx
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
