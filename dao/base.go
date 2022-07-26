package dao

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Base struct {
	rw *gorm.DB `gorm:"-"`
	ro *gorm.DB `gorm:"-"`

	ID        uint64    `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"not null; index:idx_created_at"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (m Base) DB(readOnly bool) *gorm.DB {
	if readOnly {
		return m.ro
	} else {
		return m.rw
	}
}

func (m Base) Tx(fn func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return m.rw.Transaction(fn, opts...)
}

func New(opts ...Option) Base {
	base := Base{}
	for _, opt := range opts {
		opt(&base)
	}
	return base
}
