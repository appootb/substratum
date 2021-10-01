package monitor

import (
	prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return prometheus.UnaryServerInterceptor
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return prometheus.StreamServerInterceptor
}
