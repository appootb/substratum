package mysql

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Base struct {
	tx *gorm.DB `json:"-" gorm:"-"`

	ID        uint64    `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;index:idx_created_at"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
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
