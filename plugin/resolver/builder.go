package resolver

import (
	builder "github.com/appootb/substratum/resolver"
	"google.golang.org/grpc/resolver"
)

func Init() {
	if builder.Implementor() == nil {
		builder.RegisterImplementor(&DiscoveryBuilder{})
	}
	resolver.Register(builder.Implementor())
	resolver.SetDefaultScheme(builder.Implementor().Scheme())
}

type DiscoveryBuilder struct{}

// Build creates a new resolver for the given target.
//
// gRPC dial calls Build synchronously, and fails if the returned error is
// not nil.
func (b *DiscoveryBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &DiscoveryResolver{
		target: target,
		cc:     cc,
		opts:   opts,
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}

// Scheme returns the scheme supported by this resolver.
// Scheme is defined at https://github.com/grpc/grpc/blob/master/doc/naming.md.
func (b *DiscoveryBuilder) Scheme() string {
	return "discovery"
}
