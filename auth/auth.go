package auth

import (
	"context"

	"github.com/appootb/protobuf/go/permission"
	"github.com/appootb/protobuf/go/service"
	"github.com/appootb/substratum/metadata"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	Default      = NewAlgorithmAuth()
	DefaultToken = NewJwtToken()
)

// Token interface.
type Token interface {
	// Generate a new token with specified options.
	Generate(s *permission.Secret) (string, error)
	// Refresh the token with expired time renewed.
	Refresh(s *permission.Secret) (string, error)
	// Parse and verify the token string.
	Verify(token string, level permission.TokenLevel) (*permission.Secret, error)
	// Revoke the token signed to the specified account and device.
	Revoke(accountID uint64, keyID string) error
	// Revoke all tokens signed to the specified account.
	RevokeAll(accountID uint64) error
}

type AlgorithmAuth struct {
	methods map[string]permission.TokenLevel
}

func NewAlgorithmAuth() service.Authenticator {
	return &AlgorithmAuth{
		methods: make(map[string]permission.TokenLevel),
	}
}

// Register required token level of the service.
// The map key of the parameter is the full url path of the method.
func (n *AlgorithmAuth) RegisterServiceTokenLevel(fullMethodTokenLevels map[string]permission.TokenLevel) {
	for url, level := range fullMethodTokenLevels {
		n.methods[url] = level
	}
}

// Authenticate a request specified by the full url path of the method.
func (n *AlgorithmAuth) Authenticate(ctx context.Context, fullMethod string) (*permission.Secret, error) {
	// Get request metadata.
	md := metadata.RequestMetadata(ctx)
	if md == nil {
		return nil, status.Error(codes.Unauthenticated, "request metadata not set")
	}
	if n.methods[fullMethod] == permission.TokenLevel_NONE_TOKEN && md.GetToken() == "" {
		return nil, nil
	}
	// Parse and verify account token.
	data, err := DefaultToken.Verify(md.GetToken(), n.methods[fullMethod])
	if err != nil {
		return nil, err
	}
	return data, nil
}
