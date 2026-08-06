package storage

import (
	"testing"

	"gorm.io/gorm"
)

func TestSQLConfigOptions(t *testing.T) {
	cfg := &gorm.Config{}
	WithDBDryRun()(cfg, nil)
	WithoutDBAutomaticPing()(cfg, nil)
	WithoutDBDefaultTransaction()(cfg, nil)

	if !cfg.DryRun {
		t.Fatal("DryRun not set")
	}
	if !cfg.DisableAutomaticPing {
		t.Fatal("DisableAutomaticPing not set")
	}
	if !cfg.SkipDefaultTransaction {
		t.Fatal("SkipDefaultTransaction not set")
	}
}

func TestSQLOptionsNilSafe(t *testing.T) {
	// Config-side options should ignore nil config; DB-side ignore nil db.
	WithDBDryRun()(nil, nil)
	WithoutDBAutomaticPing()(nil, nil)
	WithoutDBDefaultTransaction()(nil, nil)
	WithDBDebug()(nil, nil)
	WithDBMaxIdleConn(1)(nil, nil)
	WithDBMaxOpenConn(1)(nil, nil)
	WithDBConnMaxLifetime(0)(nil, nil)
}
