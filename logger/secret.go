package logger

import (
	"github.com/appootb/substratum/v2/proto/go/secret"
)

// RedactSecret returns log-safe account fields from authenticated context.
func RedactSecret(info *secret.Info) map[string]interface{} {
	if info == nil || info.GetAccount() == 0 {
		return nil
	}
	return map[string]interface{}{
		"account": info.GetAccount(),
		"subject": info.GetSubject().String(),
		"type":    info.GetType().String(),
		"roles":   info.GetRoles(),
	}
}
