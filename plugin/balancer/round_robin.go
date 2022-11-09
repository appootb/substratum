package balancer

import (
	"math/rand"
	"sync"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

type RoundRobinPickerBuilder struct{}

func (*RoundRobinPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	scs := make([]balancer.SubConn, 0, len(info.ReadySCs))
	for sc := range info.ReadySCs {
		scs = append(scs, sc)
	}
	return &RoundRobin{
		subConns: scs,
		// Start at a random index, as the same RR balancer rebuilds a new
		// picker when SubConn states change, and we don't want to apply excess
		// load to the first server in the list.
		next: rand.Intn(len(scs)),
	}
}

type RoundRobin struct {
	// subConns is the snapshot of the round-robin balancer when this picker was
	// created. The slice is immutable. Each Get() will do a round-robin
	// selection from it and return the selected SubConn.
	subConns []balancer.SubConn

	mu   sync.Mutex
	next int
}

func (p *RoundRobin) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	p.mu.Lock()
	sc := p.subConns[p.next]
	p.next = (p.next + 1) % len(p.subConns)
	p.mu.Unlock()
	return balancer.PickResult{SubConn: sc}, nil
}
