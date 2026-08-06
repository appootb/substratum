package configure

import (
	"strings"
	"testing"
)

func TestAddressSetMySQL(t *testing.T) {
	var addr Address
	addr.Set("mysql://user:pass@127.0.0.1:3306/app?charset=utf8mb4&parseTime=true")
	if addr.Schema != MySQL {
		t.Fatalf("schema=%s", addr.Schema)
	}
	if addr.Username != "user" || addr.Password != "pass" {
		t.Fatalf("user/pass=%s/%s", addr.Username, addr.Password)
	}
	if addr.Host != "127.0.0.1" || addr.Port != "3306" {
		t.Fatalf("host/port=%s/%s", addr.Host, addr.Port)
	}
	if addr.NameSpace != "app" {
		t.Fatalf("ns=%s", addr.NameSpace)
	}
	if addr.Params["charset"] != "utf8mb4" {
		t.Fatalf("params=%v", addr.Params)
	}
}

func TestAddressSetMultiHost(t *testing.T) {
	var addr Address
	addr.Set("rediscluster://h1:6379,h2:6379/0")
	if addr.Schema != RedisCluster {
		t.Fatalf("schema=%s", addr.Schema)
	}
	if !strings.Contains(addr.Host, ",") {
		t.Fatalf("expected multi host, got %s", addr.Host)
	}
}

func TestAddrParamsEncode(t *testing.T) {
	if AddrParams(nil).Encode("&") != "" {
		t.Fatal("nil params")
	}
	p := AddrParams{"a": "1", "b": "2"}
	enc := p.Encode("&")
	if !(enc == "a=1&b=2" || enc == "b=2&a=1") {
		t.Fatalf("encode=%s", enc)
	}
}
