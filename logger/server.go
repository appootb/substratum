package logger

import (
	"context"
	"strconv"
	"time"

	md "github.com/appootb/substratum/metadata"
	"github.com/appootb/substratum/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	LogTag    = "_BASE_."
	AccessLog = "_MSG_.access"
	ErrorLog  = "_MSG_.error"

	LogPath     = "path"
	LogConsumed = "consumed"
	LogRequest  = "request"
	LogResponse = "response"
	LogSecret   = "secret"
	LogError    = "error"
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
		consumed := time.Since(ts)
		_ = grpc.SetHeader(ctx, metadata.MD{
			LogConsumed: []string{strconv.FormatInt(consumed.Nanoseconds()/1e6, 10)},
		})

		log := Content{
			LogTag + LogPath:     info.FullMethod,
			LogTag + LogConsumed: consumed,
			LogTag + LogRequest:  req,
			LogTag + LogResponse: resp,
			LogTag + LogSecret:   service.AccountSecretFromContext(ctx),
		}
		// Access log.
		logger.Info(AccessLog, log)
		// Error log.
		if err != nil {
			log[LogTag+LogError] = err.Error()
			logger.Error(ErrorLog, log)
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
