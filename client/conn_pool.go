package client

import (
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
)

var (
	DefaultConnPool = newConnPool()
)

type ConnPool interface {
	Get(target string) *grpc.ClientConn
	Close()
}

func newConnPool() ConnPool {
	return &connPool{}
}

type connPool struct {
	sync.Map
}

func (p *connPool) Get(target string) *grpc.ClientConn {
	if cc, ok := p.Load(target); !ok {
		return cc.(*grpc.ClientConn)
	}
	cc := p.newConn(target)
	p.Store(target, cc)
	return cc
}

func (p *connPool) newConn(target string) *grpc.ClientConn {
	// TODO: support more schema
	cli, err := grpc.Dial(target,
		grpc.WithInsecure(),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  time.Second,
				MaxDelay:   time.Minute,
				Multiplier: backoff.DefaultConfig.Multiplier,
				Jitter:     backoff.DefaultConfig.Jitter,
			},
		}))
	if err != nil {
		panic("substratum: connPool dial gRPC server err:" + err.Error() + ", target:" + target)
	}
	return cli
}

func (p *connPool) Close() {
	p.Range(func(_, value interface{}) bool {
		cc := value.(*grpc.ClientConn)
		_ = cc.Close()
		return true
	})
}
