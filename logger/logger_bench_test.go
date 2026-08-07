package logger

import (
	"testing"

	"github.com/appootb/substratum/v2/proto/go/common"
	"github.com/appootb/substratum/v2/proto/go/permission"
	"github.com/appootb/substratum/v2/proto/go/secret"
)

type discardLogger struct{}

func (discardLogger) UpdateLevel(Level)                            {}
func (discardLogger) Log(Level, *common.Metadata, string, Content) {}

func BenchmarkRedactSecret(b *testing.B) {
	info := &secret.Info{
		Type:    secret.Type_CLIENT,
		Account: 42,
		KeyId:   99,
		Issuer:  "appootb",
		Subject: permission.Subject_WEB,
		Roles:   []string{"member", "editor"},
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if RedactSecret(info) == nil {
			b.Fatal("nil")
		}
	}
}

func BenchmarkRedactSecret_Nil(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if RedactSecret(nil) != nil {
			b.Fatal("expected nil")
		}
	}
}

func BenchmarkAccessLogContent(b *testing.B) {
	md := &common.Metadata{
		Product:  "app",
		Package:  "com.example",
		Version:  "1.0.0",
		Platform: common.Platform_PLATFORM_WEB,
		Locale:   "en",
	}
	info := &secret.Info{
		Type:    secret.Type_CLIENT,
		Account: 42,
		Subject: permission.Subject_WEB,
		Roles:   []string{"member"},
	}
	log := discardLogger{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := Content{
			LogPath:     "/bench.Service/Echo",
			LogConsumed: 0,
			LogRequest:  map[string]string{"k": "v"},
			LogResponse: map[string]string{"ok": "1"},
			LogSecret:   RedactSecret(info),
		}
		log.Log(InfoLevel, md, AccessLog, c)
	}
}
