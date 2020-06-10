package storage

import (
	"fmt"
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
	Port     uint16       `json:"port"`
	Database string       `json:"database"`
	Params   ConfigParams `json:"params"`
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
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		d.Username, d.Password, d.Host, d.Port, d.Database, d.Params.Encode(d.Type()))
}

type PostgreSQL struct {
	Config
}

func (d *PostgreSQL) Type() DialectType {
	return DialectPostgreSQL
}

func (d *PostgreSQL) URL() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s %s",
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
	return fmt.Sprintf("sqlserver://%s:%s@%s:%d?%s",
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
	return fmt.Sprintf("redis://:%v@%v:%v/%v", d.Password, d.Host, d.Port, d.Database)
}
