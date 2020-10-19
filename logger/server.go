package logger

import (
	"context"
	"strconv"
	"time"

	md "github.com/appootb/substratum/metadata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type loggerKey struct{}

func ContextWithLogger(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerKey{}, &Helper{
		md:     md.RequestMetadata(ctx),
		Logger: impl,
	})
}

func ContextLogger(ctx context.Context) *Helper {
	if logger := ctx.Value(loggerKey{}); logger != nil {
		return logger.(*Helper)
	}
	return &Helper{
		Logger: impl,
	}
}

// UnaryServerInterceptor returns a new unary server interceptor for access log.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ts := time.Now()
		logger := &Helper{
			Logger: impl,
			md:     md.RequestMetadata(ctx),
		}
		resp, err := handler(context.WithValue(ctx, loggerKey{}, logger), req)
		consume := time.Since(ts)
		_ = grpc.SetHeader(ctx, metadata.MD{
			"consume": []string{strconv.FormatInt(consume.Nanoseconds()/1e6, 10)},
		})

		log := Content{
			"path":     info.FullMethod,
			"consume":  consume,
			"request":  req,
			"response": resp,
		}
		// Access log.
		logger.Info("access_log", log)
		// Error log.
		if err != nil {
			log["error"] = err.Error()
			logger.Error("error_log", log)
		}
		return resp, err
	}
}

// StreamServerInterceptor returns a new streaming server interceptor for access log.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		wrapper := &ctxWrapper{stream}
		return handler(srv, wrapper)
	}
}

type ctxWrapper struct {
	grpc.ServerStream
}

func (s *ctxWrapper) Context() context.Context {
	ctx := s.ServerStream.Context()
	return ContextWithLogger(ctx)
}
