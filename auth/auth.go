package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/appootb/protobuf/go/permission"
	"github.com/appootb/protobuf/go/secret"
	"github.com/appootb/protobuf/go/service"
	"github.com/appootb/substratum/metadata"
	"github.com/appootb/substratum/util/datetime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	Default      = NewAlgorithmAuth()
	DefaultToken = NewJwtToken()
)

var (
	UnsupportedAlgorithm = errors.New("substratum: algorithm not supported")
)

// Token interface.
type Token interface {
	// Generate a new secret key.
	NewSecretKey(alg secret.Algorithm) ([]byte, error)
	// Generate a new token with specified secret info.
	Generate(s *secret.Info) (string, error)
	// Refresh the token with expired time renewed.
	Refresh(s *secret.Info) (string, error)
	// Parse the token string.
	Parse(token string) (*secret.Info, error)
}

type AlgorithmAuth struct {
	methods map[string][]permission.Subject
}

func NewAlgorithmAuth() service.Authenticator {
	return &AlgorithmAuth{
		methods: make(map[string][]permission.Subject),
	}
}

// Register required method subjects of the service.
// The map key of the parameter is the full url path of the method.
func (n *AlgorithmAuth) RegisterServiceSubjects(serviceMethodSubjects map[string][]permission.Subject) {
	for methodURL, methodSubjects := range serviceMethodSubjects {
		n.methods[methodURL] = methodSubjects
	}
}

// Authenticate a request specified by the full url path of the method.
func (n *AlgorithmAuth) Authenticate(ctx context.Context, serviceMethod string) (*secret.Info, error) {
	// Get request metadata.
	md := metadata.RequestMetadata(ctx)
	if md == nil {
		return nil, status.Error(codes.Unauthenticated, "request metadata not set")
	}
	if md.GetToken() == "" && n.IsAnonymousMethod(serviceMethod) {
		dt := time.Now().Add(-time.Minute)
		return &secret.Info{
			Roles:     []string{},
			IssuedAt:  datetime.WithTime(dt).Proto(),
			ExpiredAt: datetime.WithTime(dt).Proto(),
		}, nil
	}
	// Parse the token.
	secretInfo, err := DefaultToken.Parse(md.GetToken())
	if err != nil {
		return nil, err
	}
	// Verify the subject.
	for _, sub := range n.methods[serviceMethod] {
		if sub == secretInfo.GetSubject() {
			return secretInfo, nil
		}
	}
	return nil, status.Error(codes.PermissionDenied,
		fmt.Sprintf("subject: %v, expeced: %v", secretInfo.GetSubject(), n.methods[serviceMethod]))
}

func (n *AlgorithmAuth) IsAnonymousMethod(serviceMethod string) bool {
	for _, aud := range n.methods[serviceMethod] {
		if aud == permission.Subject_NONE {
			return true
		}
	}
	return false
}
