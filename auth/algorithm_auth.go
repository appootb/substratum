package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/appootb/protobuf/go/common"
	"github.com/appootb/protobuf/go/permission"
	"github.com/appootb/protobuf/go/secret"
	"github.com/appootb/protobuf/go/service"
	"github.com/appootb/substratum/errors"
	"github.com/appootb/substratum/metadata"
	"github.com/appootb/substratum/util/datetime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TokenParser interface {
	// Parse the token string.
	Parse(md *common.Metadata) (*secret.Info, error)
}

func NewAlgorithmAuth(client, server TokenParser) service.Authenticator {
	return &AlgorithmAuth{
		clientTokenParser: client,
		serverTokenParser: server,
		methodComponent:   make(map[string]string),
		methodSubjects:    make(map[string][]permission.Subject),
	}
}

type AlgorithmAuth struct {
	clientTokenParser TokenParser
	serverTokenParser TokenParser
	methodComponent   map[string]string
	methodSubjects    map[string][]permission.Subject
}

// Return the component name implements the service method.
func (n *AlgorithmAuth) ServiceComponentName(serviceMethod string) string {
	return n.methodComponent[serviceMethod]
}

// Register required method subjects of the service.
// The map key of the parameter is the full url path of the method.
func (n *AlgorithmAuth) RegisterServiceSubjects(component string, serviceMethodSubjects map[string][]permission.Subject) {
	for methodURL, methodSubjects := range serviceMethodSubjects {
		n.methodComponent[methodURL] = component
		n.methodSubjects[methodURL] = methodSubjects
	}
}

// Authenticate a request specified by the full url path of the method.
func (n *AlgorithmAuth) Authenticate(ctx context.Context, serviceMethod string) (*secret.Info, error) {
	dt := time.Now().Add(-time.Minute)
	anonymousMethod := n.IsAnonymousMethod(serviceMethod)
	emptySecret := &secret.Info{
		Roles:     []string{},
		IssuedAt:  datetime.WithTime(dt).Proto(),
		ExpiredAt: datetime.WithTime(dt).Proto(),
	}

	// Get request metadata.
	md := metadata.RequestMetadata(ctx)
	if md == nil {
		return nil, status.Error(codes.Unauthenticated, "request metadata not set")
	}
	if md.GetToken() == "" {
		if anonymousMethod {
			return emptySecret, nil
		}
		return nil, status.Error(codes.Unauthenticated, "token required")
	}

	// Parse the token.
	var (
		err        error
		secretInfo *secret.Info
	)
	if md.GetPlatform()&common.Platform_PLATFORM_SERVER == common.Platform_PLATFORM_SERVER {
		secretInfo, err = n.serverTokenParser.Parse(md)
	} else {
		secretInfo, err = n.clientTokenParser.Parse(md)
	}
	if err != nil {
		if anonymousMethod {
			return emptySecret, nil
		}
		switch errors.ErrorCode(err) {
		case int32(codes.AlreadyExists),
			int32(codes.FailedPrecondition),
			int32(codes.Unauthenticated):
			return nil, err
		default:
			return nil, status.Error(codes.Unauthenticated, "verify token failed")
		}
	}
	// Anonymous method
	if anonymousMethod {
		return secretInfo, nil
	}
	// Verify the subject.
	if !n.IsValidPlatform(secretInfo.GetSubject(), md.GetPlatform()) {
		return nil, status.Error(codes.Unauthenticated, "invalid token usage")
	}
	for _, sub := range n.methodSubjects[serviceMethod] {
		if (sub & secretInfo.GetSubject()) == sub {
			return secretInfo, nil
		}
	}
	return nil, status.Error(codes.PermissionDenied,
		fmt.Sprintf("subject: %v, expeced: %v", secretInfo.GetSubject(), n.methodSubjects[serviceMethod]))
}

func (n *AlgorithmAuth) IsAnonymousMethod(serviceMethod string) bool {
	for _, aud := range n.methodSubjects[serviceMethod] {
		if aud == permission.Subject_NONE {
			return true
		}
	}
	return false
}

func (n *AlgorithmAuth) IsValidPlatform(sub permission.Subject, platform common.Platform) bool {
	if platform&common.Platform_PLATFORM_SERVER > 0 {
		return (sub & permission.Subject_SERVER) == permission.Subject_SERVER
	}
	if platform&common.Platform_PLATFORM_WEB > 0 {
		return (sub & permission.Subject_WEB) == permission.Subject_WEB
	}
	if platform&common.Platform_PLATFORM_PC > 0 {
		return (sub & permission.Subject_PC) == permission.Subject_PC
	}
	if platform&common.Platform_PLATFORM_MOBILE > 0 {
		return (sub & permission.Subject_MOBILE) == permission.Subject_MOBILE
	}

	return false
}
