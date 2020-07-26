package resolver

import (
	"fmt"
	"strings"

	"github.com/appootb/substratum/discovery"
	"github.com/appootb/substratum/util/iphelper"
	"google.golang.org/grpc/resolver"
)

type DiscoveryResolver struct {
	target resolver.Target
	cc     resolver.ClientConn
	opts   resolver.BuildOptions
}

// ResolveNow will be called by gRPC to try to resolve the target name
// again. It's just a hint, resolver can ignore this if it's not necessary.
//
// It could be called multiple times concurrently.
func (r *DiscoveryResolver) ResolveNow(_ resolver.ResolveNowOptions) {
	ipPrefix := fmt.Sprintf("%v:", iphelper.LocalIP())
	nodes := discovery.Implementor().GetNodes(r.target.Endpoint)
	addrs := make([]resolver.Address, 0, len(nodes))
	for addr := range nodes {
		if strings.HasPrefix(addr, ipPrefix) || strings.HasPrefix(addr, "127.0.0.1:") {
			addrs = []resolver.Address{
				{
					Addr: addr,
				},
			}
			break
		}
		addrs = append(addrs, resolver.Address{
			Addr: addr,
		})
	}
	r.cc.UpdateState(resolver.State{
		Addresses: addrs,
	})
}

// Close closes the resolver.
func (r *DiscoveryResolver) Close() {
}
