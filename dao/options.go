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

// Option configures a Base instance.
type Option func(*Base)

// WithContext attaches ctx and resolves read/write DB handles from storage in context.
// When develop metadata is set, both handles run in Debug mode.
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

// WithDB uses tx for both read and write paths (typical inside a transaction).
func WithDB(ctx context.Context, tx *gorm.DB) Option {
	return func(base *Base) {
		base.ctx = ctx
		base.rw = tx
		base.ro = tx
	}
}

// WithPrimaryKey sets the model primary key.
func WithPrimaryKey(pk uint64) Option {
	return func(base *Base) {
		base.ID = pk
	}
}

// WithSnowflake sets CreatedAt from a snowflake ID timestamp.
func WithSnowflake(id uint64) Option {
	return func(base *Base) {
		base.CreatedAt = snowflake.Timestamp(id)
	}
}

// WithCreatedTime sets CreatedAt explicitly.
func WithCreatedTime(timestamp time.Time) Option {
	return func(base *Base) {
		base.CreatedAt = timestamp
	}
}
