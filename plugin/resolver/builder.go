package resolver

import (
	"strings"
	"sync"

	builder "github.com/appootb/substratum/v2/resolver"
	"google.golang.org/grpc/resolver"
)

const (
	DefaultSchema = "appootb"
)

func Init() {
	if builder.Implementor() == nil {
		builder.RegisterImplementor(&DiscoveryBuilder{
			serviceResolvers: map[string]*DiscoveryResolver{},
		})
	}
	resolver.Register(builder.Implementor())
	resolver.SetDefaultScheme(builder.Implementor().Scheme())
}

type DiscoveryBuilder struct {
	sync.RWMutex
	serviceResolvers map[string]*DiscoveryResolver
}

// Build creates a new resolver for the given target.
//
// gRPC dial calls Build synchronously, and fails if the returned error is
// not nil.
func (b *DiscoveryBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	_, service, _ := strings.Cut(target.URL.Path, "/")
	r := &DiscoveryResolver{
		service: service,
		cc:      cc,
		opts:    opts,
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	b.serviceResolvers[service] = r
	return r, nil
}

// Scheme returns the scheme supported by this resolver.
// Scheme is defined at https://github.com/grpc/grpc/blob/master/doc/naming.md.
func (b *DiscoveryBuilder) Scheme() string {
	return DefaultSchema
}

func (b *DiscoveryBuilder) UpdateAddresses(service string, addresses []string) error {
	b.RLock()
	defer b.RUnlock()
	r, ok := b.serviceResolvers[service]
	if !ok {
		return nil
	}
	//
	addrs := make([]resolver.Address, 0, len(addresses))
	for _, addr := range addresses {
		addrs = append(addrs, resolver.Address{
			Addr: addr,
		})
	}
	return r.cc.UpdateState(resolver.State{
		Addresses: addrs,
	})
}

func (b *DiscoveryBuilder) ReportAddressError(service string, err error) {
	b.RLock()
	defer b.RUnlock()
	r, ok := b.serviceResolvers[service]
	if ok {
		r.cc.ReportError(err)
	}
}
