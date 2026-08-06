// Package dao provides a GORM model base with request-scoped DB handles.
package dao

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// Base is embedded by application models for primary key, timestamps, and DB access.
type Base struct {
	ctx context.Context `gorm:"-"`

	rw *gorm.DB `gorm:"-"`
	ro *gorm.DB `gorm:"-"`

	ID        uint64    `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"not null; index:idx_created_at"`
	UpdatedAt time.Time `gorm:"not null"`
}

// Context returns the request context attached via options.
func (m *Base) Context() context.Context {
	return m.ctx
}

// DB returns a context-bound gorm.DB. When readOnly is true, the read replica is used if available.
func (m *Base) DB(readOnly bool) *gorm.DB {
	if readOnly {
		return m.ro.WithContext(m.ctx)
	} else {
		return m.rw.WithContext(m.ctx)
	}
}

// Tx runs fn inside a write-path transaction.
func (m *Base) Tx(fn func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return m.rw.WithContext(m.ctx).Transaction(fn, opts...)
}

// New constructs a Base applying the given options.
func New(opts ...Option) Base {
	base := Base{}
	for _, opt := range opts {
		opt(&base)
	}
	return base
}
