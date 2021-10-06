package storage

import (
	"time"

	"gorm.io/gorm"
)

type SQLOption func(*gorm.Config, *gorm.DB)

// WithDBDebug starts debug mode.
func WithDBDebug() SQLOption {
	return func(_ *gorm.Config, db *gorm.DB) {
		if db == nil {
			return
		}
		*db = *db.Debug()
	}
}

// WithDBDryRun generates sql without executing.
func WithDBDryRun() SQLOption {
	return func(cfg *gorm.Config, _ *gorm.DB) {
		if cfg == nil {
			return
		}
		cfg.DryRun = true
	}
}

// WithoutDBAutomaticPing disables automatic ping.
func WithoutDBAutomaticPing() SQLOption {
	return func(cfg *gorm.Config, _ *gorm.DB) {
		if cfg == nil {
			return
		}
		cfg.DisableAutomaticPing = true
	}
}

// WithoutDBDefaultTransaction disables single create, update, delete operations in transaction.
func WithoutDBDefaultTransaction() SQLOption {
	return func(cfg *gorm.Config, _ *gorm.DB) {
		if cfg == nil {
			return
		}
		cfg.SkipDefaultTransaction = true
	}
}

// WithDBMaxIdleConn sets the maximum number of connections in the idle connection pool.
func WithDBMaxIdleConn(maxIdleConn int) SQLOption {
	return func(_ *gorm.Config, db *gorm.DB) {
		if db == nil {
			return
		}
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.SetMaxIdleConns(maxIdleConn)
		}
	}
}

// WithDBMaxOpenConn sets the maximum number of open connections to the database.
func WithDBMaxOpenConn(maxOpenConn int) SQLOption {
	return func(_ *gorm.Config, db *gorm.DB) {
		if db == nil {
			return
		}
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.SetMaxOpenConns(maxOpenConn)
		}
	}
}

// WithDBConnMaxLifetime sets the maximum amount of time a connection may be reused.
func WithDBConnMaxLifetime(dur time.Duration) SQLOption {
	return func(_ *gorm.Config, db *gorm.DB) {
		if db == nil {
			return
		}
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.SetConnMaxLifetime(dur)
		}
	}
}
