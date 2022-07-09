package configure

import (
	"fmt"
	"net/url"
	"strings"
)

type Schema string

const (
	SQLite     Schema = "sqlite"
	MySQL      Schema = "mysql"
	PostgreSQL Schema = "postgres"
	SQLServer  Schema = "sqlserver"
	Clickhouse Schema = "clickhouse"

	Redis        Schema = "redis"
	RedisCluster Schema = "rediscluster"

	ElasticSearch6 Schema = "elasticsearch6"
	ElasticSearch7 Schema = "elasticsearch7"

	Kafka    Schema = "kafka"
	Pulsar   Schema = "pulsar"
	Pulsars  Schema = "pulsar+ssl"
	RocketMQ Schema = "rocketmq"
)

type Dialect interface {
	URL() string
}

type AddrParams map[string]string

func (m AddrParams) Encode(sep string) string {
	if m == nil || len(m) == 0 {
		return ""
	}

	var params = make([]string, 0, len(m))
	for k, v := range m {
		params = append(params, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(params, sep)
}

type Address struct {
	Schema    Schema     `json:"schema"`
	Username  string     `json:"username"`
	Password  string     `json:"password"`
	Host      string     `json:"host"`
	Port      string     `json:"port"`
	NameSpace string     `json:"namespace"`
	Params    AddrParams `json:"params"`
	RawValue  string     `json:"raw"`
}

func (c *Address) Set(addr string) {
	uri, err := url.Parse(addr)
	if err != nil {
		return
	}
	c.Schema = Schema(uri.Scheme)
	c.Username = uri.User.Username()
	c.Password, _ = uri.User.Password()
	//
	if strings.Contains(uri.Host, ",") {
		c.Host = uri.Host
	} else {
		c.Host = uri.Hostname()
		c.Port = uri.Port()
	}
	//
	c.NameSpace = strings.Trim(uri.Path, "/")
	c.RawValue = addr
	if c.Params == nil {
		c.Params = make(AddrParams, len(uri.Query()))
	}
	for k, v := range uri.Query() {
		c.Params[k] = v[0]
	}
}
