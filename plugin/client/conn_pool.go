package client

import (
	"sync"
	"time"

	"github.com/appootb/substratum/v2/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

func Init() {
	if client.Implementor() == nil {
		client.RegisterImplementor(&ConnPool{})
	}
}

type ConnPool struct {
	sync.Map
}

func (p *ConnPool) Get(target string) *grpc.ClientConn {
	if cc, ok := p.Load(target); ok {
		return cc.(*grpc.ClientConn)
	}
	cc := p.NewConn(target)
	p.Store(target, cc)
	return cc
}

func (p *ConnPool) NewConn(target string) *grpc.ClientConn {
	// TODO: support more schema
	cli, err := grpc.Dial(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
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

func (p *ConnPool) Close() {
	p.Range(func(_, value interface{}) bool {
		cc := value.(*grpc.ClientConn)
		_ = cc.Close()
		return true
	})
}
