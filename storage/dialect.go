package storage

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type DialectType string

const (
	DialectSQLite         DialectType = "sqlite3"
	DialectMySQL          DialectType = "mysql"
	DialectPostgreSQL     DialectType = "postgres"
	DialectMSSQL          DialectType = "mssql"
	DialectRedis          DialectType = "redis"
	DialectElasticSearch6 DialectType = "elasticsearch6"
	DialectElasticSearch7 DialectType = "elasticsearch7"
)

type Dialect interface {
	Type() DialectType
	Config() *Config
	URL() string
}

func ParseDialect(addr string) (Dialect, error) {
	var cfg Config
	cfg.Set(addr)
	return NewDialect(cfg)
}

func NewDialect(cfg Config) (Dialect, error) {
	if cfg.Params == nil {
		cfg.Params = map[string]string{}
	}
	switch cfg.Schema {
	case DialectSQLite:
		return &SQLite{cfg}, nil
	case DialectMySQL:
		return &MySQL{cfg}, nil
	case DialectPostgreSQL:
		return &PostgreSQL{cfg}, nil
	case DialectMSSQL:
		return &MSSQL{cfg}, nil
	case DialectRedis:
		return &Redis{cfg}, nil
	default:
		return nil, errors.New("unknown dialect type:" + string(cfg.Schema))
	}
}

type ConfigParams map[string]string

func (cp ConfigParams) Encode(dialect DialectType) string {
	if cp == nil || len(cp) == 0 {
		return ""
	}

	var (
		sep = ""
		res = make([]string, 0, len(cp))
	)

	switch dialect {
	case DialectMySQL, DialectMSSQL:
		sep = "&"
	case DialectPostgreSQL:
		sep = " "
	}
	for k, v := range cp {
		res = append(res, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(res, sep)
}

type Config struct {
	Schema   DialectType  `json:"type"`
	Username string       `json:"username"`
	Password string       `json:"password"`
	Host     string       `json:"host"`
	Port     string       `json:"port"`
	Database string       `json:"database"`
	Params   ConfigParams `json:"params"`
}

func (c *Config) Config() *Config {
	return c
}

func (c *Config) Type() DialectType {
	return c.Schema
}

func (c *Config) Set(addr string) {
	uri, err := url.Parse(addr)
	if err != nil {
		return
	}
	c.Schema = DialectType(uri.Scheme)
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

type MySQL struct {
	Config
}

func (d *MySQL) URL() string {
	d.Params["charset"] = "utf8mb4"
	d.Params["parseTime"] = "True"
	d.Params["loc"] = "Local"
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		d.Username, d.Password, d.Host, d.Port, d.Database, d.Params.Encode(d.Type()))
}

type PostgreSQL struct {
	Config
}

func (d *PostgreSQL) URL() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s %s",
		d.Host, d.Port, d.Username, d.Database, d.Password, d.Params.Encode(d.Type()))
}

type MSSQL struct {
	Config
}

func (d *MSSQL) URL() string {
	d.Params["database"] = d.Database
	return fmt.Sprintf("sqlserver://%s:%s@%s:%s?%s",
		d.Username, d.Password, d.Host, d.Port, d.Params.Encode(d.Type()))
}

type SQLite struct {
	Config
}

func (d *SQLite) URL() string {
	return d.Database
}

type Redis struct {
	Config
}

func (d *Redis) URL() string {
	return fmt.Sprintf("redis://:%s@%s:%s/%s", d.Password, d.Host, d.Port, d.Database)
}

type ElasticSearch struct {
	Config
}

func (d *ElasticSearch) URL() string {
	return fmt.Sprintf("http://%s:%s", d.Host, d.Port)
}
