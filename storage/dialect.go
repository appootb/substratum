package storage

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type DialectType string

const (
	DialectSQLite     DialectType = "sqlite3"
	DialectMySQL      DialectType = "mysql"
	DialectPostgreSQL DialectType = "postgres"
	DialectMSSQL      DialectType = "mssql"
	DialectRedis      DialectType = "redis"
)

type Dialect interface {
	Type() DialectType
	URL() string
}

func NewDialect(dialect DialectType, cfg Config, opts ...SQLOption) Dialect {
	if cfg.Params == nil {
		cfg.Params = map[string]string{}
	}
	switch dialect {
	case DialectSQLite:
		return &SQLite{cfg}
	case DialectMySQL:
		return &MySQL{cfg}
	case DialectPostgreSQL:
		return &PostgreSQL{cfg}
	case DialectMSSQL:
		return &MSSQL{cfg}
	default:
		return nil
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
	Username string       `json:"username"`
	Password string       `json:"password"`
	Host     string       `json:"host"`
	Port     string       `json:"port"`
	Database string       `json:"database"`
	Params   ConfigParams `json:"params"`
}

func (c *Config) Set(v string) {
	if err := json.Unmarshal([]byte(v), c); err == nil {
		return
	}
	uri, err := url.Parse(v)
	if err != nil {
		return
	}
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

func (d *MySQL) Type() DialectType {
	return DialectMySQL
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

func (d *PostgreSQL) Type() DialectType {
	return DialectPostgreSQL
}

func (d *PostgreSQL) URL() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s %s",
		d.Host, d.Port, d.Username, d.Database, d.Password, d.Params.Encode(d.Type()))
}

type MSSQL struct {
	Config
}

func (d *MSSQL) Type() DialectType {
	return DialectMSSQL
}

func (d *MSSQL) URL() string {
	d.Params["database"] = d.Database
	return fmt.Sprintf("sqlserver://%s:%s@%s:%s?%s",
		d.Username, d.Password, d.Host, d.Port, d.Params.Encode(d.Type()))
}

type SQLite struct {
	Config
}

func (d *SQLite) Type() DialectType {
	return DialectSQLite
}

func (d *SQLite) URL() string {
	return d.Database
}

type Redis struct {
	Config
}

func (d *Redis) Type() DialectType {
	return DialectRedis
}

func (d *Redis) URL() string {
	return fmt.Sprintf("redis://:%s@%s:%s/%s", d.Password, d.Host, d.Port, d.Database)
}
