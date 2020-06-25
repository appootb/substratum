package resolver

import (
	"sync"

	"google.golang.org/grpc/resolver"
)

var (
	Default = newBuilder()
)

var (
	once sync.Once
)

// Register the default builder.
func Register() {
	once.Do(func() {
		// Register and set the default schema.
		resolver.Register(Default)
		resolver.SetDefaultScheme(Default.Scheme())
	})
}

func newBuilder() resolver.Builder {
	return &discoveryBuilder{}
}

type discoveryBuilder struct{}

// Build creates a new resolver for the given target.
//
// gRPC dial calls Build synchronously, and fails if the returned error is
// not nil.
func (b *discoveryBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &discoveryResolver{
		target: target,
		cc:     cc,
		opts:   opts,
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}

// Scheme returns the scheme supported by this resolver.
// Scheme is defined at https://github.com/grpc/grpc/blob/master/doc/naming.md.
func (b *discoveryBuilder) Scheme() string {
	return "discovery"
}
