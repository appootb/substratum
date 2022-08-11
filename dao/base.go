package dao

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Base struct {
	ctx context.Context `gorm:"-"`

	rw *gorm.DB `gorm:"-"`
	ro *gorm.DB `gorm:"-"`

	ID        uint64    `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"not null; index:idx_created_at"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (m Base) Context() context.Context {
	return m.ctx
}

func (m Base) DB(readOnly bool) *gorm.DB {
	if readOnly {
		return m.ro.WithContext(m.ctx)
	} else {
		return m.rw.WithContext(m.ctx)
	}
}

func (m Base) Tx(fn func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return m.rw.WithContext(m.ctx).Transaction(fn, opts...)
}

func New(opts ...Option) Base {
	base := Base{}
	for _, opt := range opts {
		opt(&base)
	}
	return base
}
