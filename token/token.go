package token

import (
	"github.com/appootb/substratum/proto/go/common"
	"github.com/appootb/substratum/proto/go/secret"
)

var (
	impl Token
)

// Implementor return the token service implementor.
func Implementor() Token {
	return impl
}

// RegisterImplementor registers the token service implementor.
func RegisterImplementor(token Token) {
	impl = token
}

// Token interface.
type Token interface {
	// NewSecretKey creates a new secret key.
	NewSecretKey(alg secret.Algorithm) ([]byte, error)
	// Generate a new token with specified secret info.
	Generate(s *secret.Info) (string, error)
	// Refresh the token with expired time renewed.
	Refresh(s *secret.Info) (string, error)
	// Parse the metadata.
	Parse(md *common.Metadata) (*secret.Info, error)
	// ParseRaw parses a token string.
	ParseRaw(token string) (*secret.Info, error)
}
