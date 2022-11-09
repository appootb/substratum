package balancer

import (
	builder "github.com/appootb/substratum/v2/balancer"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/balancer/roundrobin"
)

const Name = "round_robin"

func Init() {
	if builder.Implementor() == nil {
		impl := base.NewBalancerBuilder(roundrobin.Name, &RoundRobinPickerBuilder{}, base.Config{HealthCheck: true})
		builder.RegisterImplementor(impl)
	}
	balancer.Register(builder.Implementor())
}
