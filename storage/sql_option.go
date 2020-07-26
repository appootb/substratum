package storage

import (
	"time"

	"github.com/jinzhu/gorm"
)

type SQLOption func(*gorm.DB)

func WithDBDump() SQLOption {
	return func(db *gorm.DB) {
		db.LogMode(true)
	}
}

func WithDBMaxIdleConn(maxIdleConn int) SQLOption {
	return func(db *gorm.DB) {
		db.DB().SetMaxIdleConns(maxIdleConn)
	}
}

func WithDBMaxOpenConn(maxOpenConn int) SQLOption {
	return func(db *gorm.DB) {
		db.DB().SetMaxOpenConns(maxOpenConn)
	}
}

func WithDBConnMaxLifetime(dur time.Duration) SQLOption {
	return func(db *gorm.DB) {
		db.DB().SetConnMaxLifetime(dur)
	}
}
