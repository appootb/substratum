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
	LogTag        = "_BASE_."
	AccessLog     = "_MSG_.access"
	ErrorLog      = "_MSG_.error"
	UpstreamLog   = "_MSG_.upstream"
	DownstreamLog = "_MSG_.downstream"
	StreamingLog  = "_MSG_.streaming"

	Consumed      = "consumed"
	LogConsumed   = LogTag + Consumed
	LogPath       = LogTag + "path"
	LogRequest    = LogTag + "request"
	LogResponse   = LogTag + "response"
	LogUpstream   = LogTag + "upstream"
	LogDownstream = LogTag + "downstream"
	LogSecret     = LogTag + "secret"
	LogError      = LogTag + "error"
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
		var (
			ts     = time.Now()
			logger = newHelper(ctx)
		)
		resp, err := handler(context.WithValue(ctx, loggerKey{}, logger), req)
		consumed := time.Since(ts)
		_ = grpc.SetHeader(ctx, metadata.MD{
			Consumed: []string{strconv.FormatInt(consumed.Nanoseconds()/1e6, 10)},
		})

		log := Content{
			LogPath:     info.FullMethod,
			LogConsumed: consumed,
			LogRequest:  req,
			LogResponse: resp,
			LogSecret:   service.AccountSecretFromContext(ctx),
		}
		if err != nil {
			// Error log.
			log[LogError] = err.Error()
			logger.Error(ErrorLog, log)
		} else {
			// Access log.
			logger.Info(AccessLog, log)
		}
		return resp, err
	}
}

// StreamServerInterceptor returns a new streaming server interceptor for access log.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		logger := newHelper(stream.Context())
		err := handler(srv, &ctxWrapper{
			ServerStream: stream,
			info:         info,
			logger:       logger,
		})
		log := Content{
			LogPath:   info.FullMethod,
			LogSecret: service.AccountSecretFromContext(stream.Context()),
		}
		if err != nil {
			// Error log.
			log[LogError] = err.Error()
			logger.Error(ErrorLog, log)
		} else {
			// Streaming log.
			logger.Info(StreamingLog, log)
		}
		return err
	}
}

type ctxWrapper struct {
	grpc.ServerStream
	info   *grpc.StreamServerInfo
	logger *Helper
}

func (s *ctxWrapper) Context() context.Context {
	ctx := s.ServerStream.Context()
	return context.WithValue(ctx, loggerKey{}, s.logger)
}

func (s *ctxWrapper) SendMsg(m interface{}) error {
	err := s.ServerStream.SendMsg(m)
	log := Content{
		LogPath:       s.info.FullMethod,
		LogDownstream: m,
		LogSecret:     service.AccountSecretFromContext(s.ServerStream.Context()),
	}
	if err != nil {
		// Error log.
		log[LogError] = err.Error()
		s.logger.Error(ErrorLog, log)
	} else {
		// Downstream log.
		s.logger.Info(DownstreamLog, log)
	}
	return err
}

func (s *ctxWrapper) RecvMsg(m interface{}) error {
	err := s.ServerStream.RecvMsg(m)
	log := Content{
		LogPath:     s.info.FullMethod,
		LogUpstream: m,
		LogSecret:   service.AccountSecretFromContext(s.ServerStream.Context()),
	}
	if err != nil {
		// Error log.
		log[LogError] = err.Error()
		s.logger.Error(ErrorLog, log)
	} else {
		// Upstream log.
		s.logger.Info(UpstreamLog, log)
	}
	return err
}
