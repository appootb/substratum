package server

import (
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/appootb/substratum/gateway"
	"github.com/appootb/substratum/logger"
	"github.com/appootb/substratum/rpc"
	"github.com/appootb/substratum/util/ip"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

type ServeMux struct {
	connAddr        string
	rpcListener     net.Listener
	gatewayListener net.Listener

	rpcSrv     *grpc.Server
	httpMux    *http.ServeMux
	gatewayMux *runtime.ServeMux
}

func NewServeMux(rpcPort, gatewayPort uint16) (*ServeMux, error) {
	var err error
	m := &ServeMux{
		rpcSrv:     rpc.New(rpc.DefaultOptions),
		httpMux:    http.NewServeMux(),
		gatewayMux: gateway.New(gateway.DefaultOptions),
	}
	m.connAddr = fmt.Sprintf("%s:%d", m.ConnIP(), rpcPort)
	m.rpcListener, err = net.Listen("tcp", fmt.Sprintf(":%d", rpcPort))
	if err != nil {
		return nil, err
	}
	m.gatewayListener, err = net.Listen("tcp", fmt.Sprintf(":%d", gatewayPort))
	if err != nil {
		return nil, err
	}
	m.httpMux.Handle("/", m.gatewayMux)
	return m, nil
}

func (m *ServeMux) RPCServer() *grpc.Server {
	return m.rpcSrv
}

func (m *ServeMux) GatewayMux() *runtime.ServeMux {
	return m.gatewayMux
}

func (m *ServeMux) Handle(pattern string, handler http.Handler) {
	m.httpMux.Handle(pattern, handler)
}

func (m *ServeMux) HandleFunc(pattern string, handler http.HandlerFunc) {
	m.httpMux.HandleFunc(pattern, handler)
}

func (m *ServeMux) Serve() {
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

func (m *ServeMux) ConnIP() string {
	ipv4, err := ip.LocalIP(true)
	if err != nil {
		logger.Error("ip.LocalIP v4", logger.Content{
			"err": err.Error(),
		})
	}
	if ipv4 != nil {
		return ipv4.String()
	}
	ipv6, err := ip.LocalIP(false)
	if err != nil {
		logger.Error("ip.LocalIP v6", logger.Content{
			"err": err.Error(),
		})
	}
	if ipv6 != nil {
		return ipv6.String()
	}
	return "127.0.0.1"
}

func (m *ServeMux) ConnAddr() string {
	return m.connAddr
}
