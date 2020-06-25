package storage

import (
	"reflect"
	"testing"
)

func Test_ParseRedisUrlConfig(t *testing.T) {
	cfg := &Redis{
		Config: Config{
			Schema:   DialectRedis,
			Password: "pwd",
			Host:     "127.0.0.1",
			Port:     "6379",
			Database: "1",
			Params:   ConfigParams{},
		},
	}
	addr := "redis://:pwd@127.0.0.1:6379/1"
	dialect, err := ParseDialect(addr)
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(dialect, cfg) {
		t.Fatal("expect", cfg, "actual", dialect)
	}
}

func Test_ParseMySQLUrlConfig(t *testing.T) {
	cfg := &MySQL{
		Config: Config{
			Schema:   DialectMySQL,
			Username: "root",
			Password: "pwd",
			Host:     "127.0.0.1",
			Port:     "3306",
			Database: "test",
			Params: ConfigParams{
				"charset":   "utf8mb4",
				"parseTime": "True",
			},
		},
	}
	addr := "mysql://root:pwd@127.0.0.1:3306/test?charset=utf8mb4&parseTime=True"
	dialect, err := ParseDialect(addr)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(dialect, cfg) {
		t.Fatal("expect", cfg, "actual", dialect)
	}
}
