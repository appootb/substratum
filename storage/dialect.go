package storage

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"gorm.io/gorm"
)

var (
	sqlImpl SQLDialect
)

// SQLDialectImplementor returns the gorm sql Dialector implementor.
func SQLDialectImplementor() SQLDialect {
	return sqlImpl
}

// RegisterSQLDialectImplementor registers gorm sql Dialector.
func RegisterSQLDialectImplementor(sql SQLDialect) {
	sqlImpl = sql
}

type SQLDialect interface {
	Open(Config) gorm.Dialector
}

type Schema string

const (
	SQLite         Schema = "sqlite"
	MySQL          Schema = "mysql"
	PostgreSQL     Schema = "postgres"
	SQLServer      Schema = "sqlserver"
	Clickhouse     Schema = "clickhouse"
	Redis          Schema = "redis"
	RedisCluster   Schema = "rediscluster"
	ElasticSearch6 Schema = "elasticsearch6"
	ElasticSearch7 Schema = "elasticsearch7"
)

type Dialect interface {
	URL() string
}

type ConfigParams map[string]string

func (cp ConfigParams) Encode(schema Schema) string {
	if cp == nil || len(cp) == 0 {
		return ""
	}

	var (
		sep = "&"
		res = make([]string, 0, len(cp))
	)

	if schema == PostgreSQL {
		sep = " "
	}
	for k, v := range cp {
		res = append(res, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(res, sep)
}

type Config struct {
	Schema   Schema       `json:"type"`
	Username string       `json:"username"`
	Password string       `json:"password"`
	Host     string       `json:"host"`
	Port     string       `json:"port"`
	Database string       `json:"database"`
	Params   ConfigParams `json:"params"`
}

func (c *Config) Dialect() (Dialect, error) {
	if c.Params == nil {
		c.Params = map[string]string{}
	}
	switch c.Schema {
	case Redis, RedisCluster:
		return &redisCache{c}, nil
	case ElasticSearch6, ElasticSearch7:
		return &elasticSearch{c}, nil
	default:
		return nil, errors.New("unknown internal dialect:" + string(c.Schema))
	}
}

func (c *Config) Set(addr string) {
	uri, err := url.Parse(addr)
	if err != nil {
		return
	}
	c.Schema = Schema(uri.Scheme)
	c.Username = uri.User.Username()
	c.Password, _ = uri.User.Password()
	c.Host = uri.Hostname()
	c.Port = uri.Port()
	c.Database = strings.Trim(uri.Path, "/")
	if c.Params == nil {
		c.Params = make(ConfigParams, len(uri.Query()))
	}
	for k, v := range uri.Query() {
		c.Params[k] = v[0]
	}
}

type redisCache struct {
	*Config
}

func (d *redisCache) URL() string {
	return fmt.Sprintf("%s://:%s@%s:%s/%s", Redis, d.Password, d.Host, d.Port, d.Database)
}

type elasticSearch struct {
	*Config
}

func (d *elasticSearch) URL() string {
	return fmt.Sprintf("http://%s:%s", d.Host, d.Port)
}
