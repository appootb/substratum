package resolver

import (
	"github.com/appootb/substratum/discovery"
	"google.golang.org/grpc/resolver"
)

type discoveryResolver struct {
	target resolver.Target
	cc     resolver.ClientConn
	opts   resolver.BuildOptions
}

// ResolveNow will be called by gRPC to try to resolve the target name
// again. It's just a hint, resolver can ignore this if it's not necessary.
//
// It could be called multiple times concurrently.
func (r *discoveryResolver) ResolveNow(_ resolver.ResolveNowOptions) {
	nodes := discovery.DefaultService.GetNodes(r.target.Endpoint)
	addrs := make([]resolver.Address, 0, len(nodes))
	for addr := range nodes {
		addrs = append(addrs, resolver.Address{
			Addr: addr,
		})
	}
	r.cc.UpdateState(resolver.State{
		Addresses: addrs,
	})
}

// Close closes the resolver.
func (r *discoveryResolver) Close() {
}
