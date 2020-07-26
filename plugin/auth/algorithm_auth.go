package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/appootb/protobuf/go/common"
	"github.com/appootb/protobuf/go/permission"
	"github.com/appootb/protobuf/go/secret"
	"github.com/appootb/substratum/auth"
	"github.com/appootb/substratum/metadata"
	"github.com/appootb/substratum/token"
	"github.com/appootb/substratum/util/datetime"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Init() {
	if auth.Implementor() == nil {
		auth.RegisterImplementor(&AlgorithmAuth{
			methodComponent: make(map[string]string),
			methodSubjects:  make(map[string][]permission.Subject),
		})
	}
}

type AlgorithmAuth struct {
	methodComponent map[string]string
	methodSubjects  map[string][]permission.Subject
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
			md.Account = proto.Uint64(0)
			return emptySecret, nil
		}
		return nil, status.Error(codes.Unauthenticated, "token required")
	}
	// Parse the token.
	secretInfo, err := token.Implementor().Parse(md.GetToken())
	if err != nil {
		if anonymousMethod {
			md.Account = proto.Uint64(0)
			return emptySecret, nil
		}
		return nil, err
	}
	// Anonymous method
	if anonymousMethod {
		md.Account = &secretInfo.Account
		return secretInfo, nil
	}
	// Verify account and secret info
	if md.GetAccount() != secretInfo.Account {
		if md.GetIsDebug() {
			md.Account = &secretInfo.Account
		} else {
			return nil, status.Error(codes.Unauthenticated, "invalid account token")
		}
	}
	// Verify the subject.
	if !n.IsValidPlatform(secretInfo.GetSubject(), md.GetPlatform()) {
		return nil, status.Error(codes.Unauthenticated, "invalid token usage")
	}
	for _, sub := range n.methodSubjects[serviceMethod] {
		if sub == secretInfo.GetSubject() {
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
	switch platform {
	case common.Platform_PLATFORM_H5, common.Platform_PLATFORM_WEB, common.Platform_PLATFORM_CHROME:
		return sub == permission.Subject_WEB
	case common.Platform_PLATFORM_LINUX, common.Platform_PLATFORM_WINDOWS, common.Platform_PLATFORM_DARWIN:
		return sub == permission.Subject_PC
	case common.Platform_PLATFORM_ANDROID, common.Platform_PLATFORM_IOS:
		return sub == permission.Subject_MOBILE
	case common.Platform_PLATFORM_SERVER:
		return sub == permission.Subject_SERVER
	default:
		return false
	}
}
