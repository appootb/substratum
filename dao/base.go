package dao

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Base struct {
	tx *gorm.DB `gorm:"-"`

	ID        uint64    `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"not null; index:idx_created_at"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (m Base) DB() *gorm.DB {
	return m.tx
}

func New(opts ...Option) Base {
	base := Base{}
	for _, opt := range opts {
		opt(&base)
	}
	return base
}
