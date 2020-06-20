package storage

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Test_SetJsonConfig(t *testing.T) {
	var c Config
	cfg := &Config{
		Username: "Test_SetJsonConfig",
		Password: "pwd",
		Host:     "127.0.0.1",
		Port:     "3306",
		Database: "Test_SetJsonConfig",
		Params:   ConfigParams{},
	}
	v, err := json.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}
	c.Set(string(v))
	if !reflect.DeepEqual(&c, cfg) {
		t.Fatal("expect", cfg, "actual", c)
	}
}

func Test_SetRedisUrlConfig(t *testing.T) {
	var c Config
	cfg := &Config{
		Password: "pwd",
		Host:     "127.0.0.1",
		Port:     "6379",
		Database: "1",
		Params:   ConfigParams{},
	}
	uri := "redis://:pwd@127.0.0.1:6379/1"
	c.Set(uri)
	if !reflect.DeepEqual(&c, cfg) {
		t.Fatal("expect", cfg, "actual", c)
	}
}

func Test_SetMySQLUrlConfig(t *testing.T) {
	var c Config
	cfg := &Config{
		Username: "root",
		Password: "pwd",
		Host:     "127.0.0.1",
		Port:     "3306",
		Database: "test",
		Params: ConfigParams{
			"charset":   "utf8mb4",
			"parseTime": "True",
		},
	}
	uri := "mysql://root:pwd@127.0.0.1:3306/test?charset=utf8mb4&parseTime=True"
	c.Set(uri)
	if !reflect.DeepEqual(&c, cfg) {
		t.Fatal("expect", cfg, "actual", c)
	}
}
