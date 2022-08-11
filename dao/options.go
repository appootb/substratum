package dao

import (
	"context"
	"time"

	"github.com/appootb/substratum/v2/metadata"
	"github.com/appootb/substratum/v2/service"
	"github.com/appootb/substratum/v2/storage"
	"github.com/appootb/substratum/v2/util/snowflake"
	"gorm.io/gorm"
)

type Option func(*Base)

func WithContext(ctx context.Context) Option {
	return func(base *Base) {
		base.ctx = ctx
		component := service.ComponentNameFromContext(ctx)
		base.rw = storage.ContextStorage(ctx, component).GetDB()
		base.ro = storage.ContextStorage(ctx, component).GetDB(true)
		if md := metadata.IncomingMetadata(ctx); md != nil && md.GetIsDevelop() {
			base.rw = base.rw.Debug()
			base.ro = base.ro.Debug()
		}
	}
}

func WithDB(ctx context.Context, tx *gorm.DB) Option {
	return func(base *Base) {
		base.ctx = ctx
		base.rw = tx
		base.ro = tx
	}
}

func WithPrimaryKey(pk uint64) Option {
	return func(base *Base) {
		base.ID = pk
	}
}

func WithSnowflake(id uint64) Option {
	return func(base *Base) {
		base.CreatedAt = snowflake.Timestamp(id)
	}
}

func WithCreatedTime(timestamp time.Time) Option {
	return func(base *Base) {
		base.CreatedAt = timestamp
	}
}
