package mysql

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
)

type Option func(*Base)

func WithContext(ctx context.Context) Option {
	return func(base *Base) {
		base.ctx = ctx
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
