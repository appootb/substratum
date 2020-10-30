package token

import (
	"github.com/appootb/protobuf/go/common"
	"github.com/appootb/protobuf/go/secret"
)

var (
	impl Token
)

// Return the service implementor.
func Implementor() Token {
	return impl
}

// Register service implementor.
func RegisterImplementor(token Token) {
	impl = token
}

// Token interface.
type Token interface {
	// Generate a new secret key.
	NewSecretKey(alg secret.Algorithm) ([]byte, error)
	// Generate a new token with specified secret info.
	Generate(s *secret.Info) (string, error)
	// Refresh the token with expired time renewed.
	Refresh(s *secret.Info) (string, error)
	// Parse the metadata.
	Parse(md *common.Metadata) (*secret.Info, error)
	// Parse raw token string.
	ParseRaw(token string) (*secret.Info, error)
}
