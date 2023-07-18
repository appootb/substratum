package metadata

import (
	"context"

	"github.com/appootb/substratum/proto/go/common"
	"google.golang.org/grpc"
)

type metadataKey struct{}

// UnaryServerInterceptor returns a new unary server interceptor for parsing request metadata.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md := ParseIncomingMetadata(ctx)
		return handler(context.WithValue(ctx, metadataKey{}, md), req)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor for parsing request metadata.
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
	md := ParseIncomingMetadata(ctx)
	return context.WithValue(ctx, metadataKey{}, md)
}

func AppendOutgoingMetadata(ctx context.Context, md *common.Metadata) context.Context {
	if md == nil {
		return context.WithValue(ctx, metadataKey{}, &common.Metadata{})
	}
	oldMD := IncomingMetadata(ctx)
	if oldMD == nil {
		return context.WithValue(ctx, metadataKey{}, md)
	}
	//
	if md.GetPlatform() > 0 {
		oldMD.Platform = md.Platform
	}
	if md.GetNetwork() > 0 {
		oldMD.Network = md.Network
	}
	if md.GetPackage() != "" {
		oldMD.Package = md.Package
	}
	if md.GetVersion() != "" {
		oldMD.Version = md.Version
	}
	if md.GetOsVersion() != "" {
		oldMD.OsVersion = md.OsVersion
	}
	if md.GetBrand() != "" {
		oldMD.Brand = md.Brand
	}
	if md.GetModel() != "" {
		oldMD.Model = md.Model
	}
	if md.GetDeviceId() != "" {
		oldMD.DeviceId = md.DeviceId
	}
	if md.GetIsEmulator() {
		oldMD.IsEmulator = md.IsEmulator
	}
	if md.GetIsDebug() {
		oldMD.IsDebug = md.IsDebug
	}
	if md.GetClientIp() != "" {
		oldMD.ClientIp = md.ClientIp
	}
	if md.GetChannel() != "" {
		oldMD.Channel = md.Channel
	}
	if md.GetProduct() != "" {
		oldMD.Product = md.Product
	}
	if md.GetTraceId() != "" {
		oldMD.TraceId = md.TraceId
	}
	if md.GetRiskId() != "" {
		oldMD.RiskId = md.RiskId
	}
	if md.GetUuid() != "" {
		oldMD.Uuid = md.Uuid
	}
	if md.GetUdid() != "" {
		oldMD.Udid = md.Udid
	}
	if md.GetUserAgent() != "" {
		oldMD.UserAgent = md.UserAgent
	}
	if md.GetDeviceMac() != "" {
		oldMD.DeviceMac = md.DeviceMac
	}
	if md.GetAndroidId() != "" {
		oldMD.AndroidId = md.AndroidId
	}
	return context.WithValue(ctx, metadataKey{}, oldMD)
}

func ContextWithProduct(ctx context.Context, product string) context.Context {
	return context.WithValue(ctx, metadataKey{}, &common.Metadata{
		Product: product,
	})
}

func IncomingMetadata(ctx context.Context) *common.Metadata {
	if md := ctx.Value(metadataKey{}); md != nil {
		return md.(*common.Metadata)
	}
	return nil
}
