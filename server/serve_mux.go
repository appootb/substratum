package server

import (
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/appootb/substratum/v2/gateway"
	"github.com/appootb/substratum/v2/logger"
	"github.com/appootb/substratum/v2/rpc"
	"github.com/appootb/substratum/v2/service"
	"github.com/appootb/substratum/v2/util/iphelper"
	prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

type ServeMux struct {
	metrics         bool
	connAddr        string
	rpcListener     net.Listener
	gatewayListener net.Listener

	rpcSrv     *grpc.Server
	httpMux    *http.ServeMux
	gatewayMux *runtime.ServeMux
}

func NewServeMux(rpcPort, gatewayPort uint16, metrics bool) (*ServeMux, error) {
	var err error
	m := &ServeMux{
		rpcSrv: rpc.New(
			rpc.NewOptions(
				rpc.WithDefaultKeepaliveOption(),
				rpc.WithDefaultUnaryInterceptors(),
				rpc.WithDefaultStreamInterceptors(),
			),
		),
		metrics:    metrics,
		httpMux:    http.NewServeMux(),
		gatewayMux: gateway.New(gateway.DefaultOptions),
	}
	m.connAddr = fmt.Sprintf("%s:%d", iphelper.LocalIP(), rpcPort)
	m.rpcListener, err = net.Listen("tcp", fmt.Sprintf(":%d", rpcPort))
	if err != nil {
		return nil, err
	}
	m.gatewayListener, err = net.Listen("tcp", fmt.Sprintf(":%d", gatewayPort))
	if err != nil {
		return nil, err
	}
	m.httpMux.Handle("/", m.gatewayMux)
	if metrics {
		m.httpMux.Handle("/metrics", promhttp.Handler())
	}
	return m, nil
}

func (m *ServeMux) RPCServer() *grpc.Server {
	return m.rpcSrv
}

func (m *ServeMux) HTTPMux(comp string) service.HttpHandler {
	return &httpServeMux{
		component: comp,
		serveMux:  m.httpMux,
	}
}

func (m *ServeMux) GatewayMux() *runtime.ServeMux {
	return m.gatewayMux
}

func (m *ServeMux) Serve() {
	if m.metrics {
		prometheus.Register(m.rpcSrv)
	}
	//
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		wg.Done()
		err := m.rpcSrv.Serve(m.rpcListener)
		if err != nil {
			logger.Error("rpc_server", logger.Content{
				"server": "gRPC",
				"addr":   m.rpcListener.Addr(),
				"err":    err.Error(),
			})
		}
	}()
	go func() {
		wg.Done()
		err := http.Serve(m.gatewayListener, m.httpMux)
		if err != nil {
			logger.Error("gateway_server", logger.Content{
				"server": "gateway",
				"addr":   m.gatewayListener.Addr(),
				"err":    err.Error(),
			})
		}
	}()

	wg.Wait()
}

func (m *ServeMux) ConnAddr() string {
	return m.connAddr
}
