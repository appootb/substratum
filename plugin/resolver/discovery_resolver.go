package resolver

import (
	"github.com/appootb/substratum/v2/discovery"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
)

type DiscoveryResolver struct {
	service string

	cc   resolver.ClientConn
	opts resolver.BuildOptions
}

// ResolveNow will be called by gRPC to try to resolve the target name
// again. It's just a hint, resolver can ignore this if it's not necessary.
//
// It could be called multiple times concurrently.
func (r *DiscoveryResolver) ResolveNow(_ resolver.ResolveNowOptions) {
	addresses := discovery.Implementor().GetAddresses(r.service)
	if len(addresses) == 0 {
		r.cc.ReportError(status.Error(codes.Unavailable, r.service))
		return
	}
	//
	_ = r.cc.UpdateState(resolver.State{
		Addresses: addresses,
	})
}

// Close closes the resolver.
func (r *DiscoveryResolver) Close() {
}
