package rpc

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/appootb/substratum/v2/auth"
	sctx "github.com/appootb/substratum/v2/context"
	"github.com/appootb/substratum/v2/errors"
	"github.com/appootb/substratum/v2/logger"
	mdpkg "github.com/appootb/substratum/v2/metadata"
	"github.com/appootb/substratum/v2/monitor"
	"github.com/appootb/substratum/v2/plugin/credential"
	plugintoken "github.com/appootb/substratum/v2/plugin/token"
	"github.com/appootb/substratum/v2/proto/go/common"
	"github.com/appootb/substratum/v2/proto/go/permission"
	"github.com/appootb/substratum/v2/proto/go/secret"
	"github.com/appootb/substratum/v2/recovery"
	"github.com/appootb/substratum/v2/service"
	"github.com/appootb/substratum/v2/token"
	"github.com/appootb/substratum/v2/util/hash"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Framework interceptor-chain microbenchmarks.
//
// Run:
//
//	go test ./rpc/ -bench='Benchmark' -benchmem -count=3
//
// Categories:
//   - Baseline / chain composition cost
//   - Per-interceptor cost (isolated)
//   - Default 8-interceptor stack (anonymous + JWT)
//   - In-process gRPC (bufconn): no interceptor vs full stack

const (
	benchMethodAnonymous = "/bench.Service/Anonymous"
	benchMethodAuthed    = "/bench.Service/Authed"
	benchBufSize         = 1 << 20
)

type nopLogger struct{}

func (nopLogger) UpdateLevel(logger.Level) {}
func (nopLogger) Log(logger.Level, *common.Metadata, string, logger.Content) {
}

type emptyPrompter struct{}

func (emptyPrompter) Translate(string, int32) string { return "" }

var (
	benchOnce sync.Once
	benchAuth service.Authenticator
	benchJWT  string
)

func ensureBenchDeps(b *testing.B) {
	b.Helper()
	benchOnce.Do(func() {
		if logger.Implementor() == nil {
			logger.RegisterImplementor(nopLogger{})
		}
		if errors.Implementor() == nil {
			errors.RegisterImplementor(emptyPrompter{})
		}
		credential.Init()
		plugintoken.Init()

		j := token.Implementor()
		benchAuth = auth.NewAlgorithmAuth(j, j)
		benchAuth.RegisterServiceSubjects("bench", map[string][]permission.Subject{
			benchMethodAnonymous: {permission.Subject_NONE},
			benchMethodAuthed:    {permission.Subject_SERVER},
		}, nil)

		now := time.Now()
		raw, err := j.Generate(&secret.Info{
			Type:      secret.Type_SERVER,
			Algorithm: secret.Algorithm_HMAC,
			Issuer:    "bench",
			Account:   10001,
			KeyId:     hash.Sum("rpc-framework-bench"),
			Subject:   permission.Subject_SERVER,
			Roles:     []string{"member"},
			IssuedAt:  timestamppb.New(now),
			ExpiredAt: timestamppb.New(now.Add(time.Hour)),
		})
		if err != nil {
			panic(err)
		}
		benchJWT = raw
	})
}

func identityInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(ctx, req)
	}
}

func benchHandler(ctx context.Context, req interface{}) (interface{}, error) {
	return req, nil
}

func benchInfo(method string) *grpc.UnaryServerInfo {
	return &grpc.UnaryServerInfo{FullMethod: method}
}

func incomingMD(pairs ...string) context.Context {
	return metadata.NewIncomingContext(context.Background(), metadata.Pairs(pairs...))
}

func fullMetadataCtx(token string) context.Context {
	pairs := []string{
		mdpkg.KeyProduct, "bench",
		mdpkg.KeyPackage, "github.com/appootb/substratum",
		mdpkg.KeyVersion, "2.0.0",
		mdpkg.KeyPlatform, fmt.Sprintf("%d", int(common.Platform_PLATFORM_SERVER)),
		mdpkg.KeyNetwork, "1",
		mdpkg.KeyTimestamp, fmt.Sprintf("%d", time.Now().Unix()),
		mdpkg.KeyLocale, "en",
		mdpkg.KeyUserAgent, "bench-agent",
		mdpkg.KeyOriginalIP, "127.0.0.1",
	}
	if token != "" {
		pairs = append(pairs, mdpkg.KeyToken, token)
	}
	return incomingMD(pairs...)
}

// defaultFrameworkInterceptors mirrors BuildUnaryChain without phase slots,
// so benches stay independent of the process-wide frozen chain.
func defaultFrameworkInterceptors(authenticator service.Authenticator) []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		recovery.UnaryServerInterceptor(),
		mdpkg.UnaryServerInterceptor(),
		errors.UnaryResponseInterceptor(),
		monitor.UnaryServerInterceptor(),
		service.UnaryServerInterceptor(authenticator),
		logger.UnaryServerInterceptor(),
		validator.UnaryServerInterceptor(),
		sctx.UnaryServerInterceptor(),
	}
}

func runUnary(b *testing.B, chain grpc.UnaryServerInterceptor, ctx context.Context, info *grpc.UnaryServerInfo, req interface{}) {
	b.Helper()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := chain(ctx, req, info, benchHandler); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHandlerOnly(b *testing.B) {
	ctx := context.Background()
	req := &emptypb.Empty{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := benchHandler(ctx, req); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkChainIdentity(b *testing.B) {
	for _, n := range []int{1, 4, 8, 12, 16} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			interceptors := make([]grpc.UnaryServerInterceptor, n)
			for i := range interceptors {
				interceptors[i] = identityInterceptor()
			}
			chain := chainUnaryServer(interceptors...)
			runUnary(b, chain, context.Background(), benchInfo(benchMethodAnonymous), &emptypb.Empty{})
		})
	}
}

func BenchmarkInterceptor_Recovery(b *testing.B) {
	ensureBenchDeps(b)
	runUnary(b, recovery.UnaryServerInterceptor(), context.Background(), benchInfo(benchMethodAnonymous), &emptypb.Empty{})
}

func BenchmarkInterceptor_Metadata(b *testing.B) {
	ensureBenchDeps(b)
	ctx := fullMetadataCtx("")
	runUnary(b, mdpkg.UnaryServerInterceptor(), ctx, benchInfo(benchMethodAnonymous), &emptypb.Empty{})
}

func BenchmarkInterceptor_Errors(b *testing.B) {
	ensureBenchDeps(b)
	// errors interceptor reads IncomingMetadata via the metadata package key.
	chained := chainUnaryServer(mdpkg.UnaryServerInterceptor(), errors.UnaryResponseInterceptor())
	runUnary(b, chained, fullMetadataCtx(""), benchInfo(benchMethodAnonymous), &emptypb.Empty{})
}

func BenchmarkInterceptor_Monitor(b *testing.B) {
	ensureBenchDeps(b)
	runUnary(b, monitor.UnaryServerInterceptor(), context.Background(), benchInfo(benchMethodAnonymous), &emptypb.Empty{})
}

func BenchmarkInterceptor_AuthAnonymous(b *testing.B) {
	ensureBenchDeps(b)
	// Auth requires metadata interceptor key + token handling.
	chain := chainUnaryServer(mdpkg.UnaryServerInterceptor(), service.UnaryServerInterceptor(benchAuth))
	runUnary(b, chain, fullMetadataCtx(""), benchInfo(benchMethodAnonymous), &emptypb.Empty{})
}

func BenchmarkInterceptor_AuthJWT(b *testing.B) {
	ensureBenchDeps(b)
	chain := chainUnaryServer(mdpkg.UnaryServerInterceptor(), service.UnaryServerInterceptor(benchAuth))
	runUnary(b, chain, fullMetadataCtx(benchJWT), benchInfo(benchMethodAuthed), &emptypb.Empty{})
}

func BenchmarkInterceptor_Logger(b *testing.B) {
	ensureBenchDeps(b)
	// Logger clears token on metadata; metadata interceptor must run first.
	chain := chainUnaryServer(mdpkg.UnaryServerInterceptor(), logger.UnaryServerInterceptor())
	runUnary(b, chain, fullMetadataCtx(""), benchInfo(benchMethodAnonymous), &emptypb.Empty{})
}

func BenchmarkInterceptor_Validator(b *testing.B) {
	ensureBenchDeps(b)
	runUnary(b, validator.UnaryServerInterceptor(), context.Background(), benchInfo(benchMethodAnonymous), &emptypb.Empty{})
}

func BenchmarkInterceptor_ContextInjectors(b *testing.B) {
	ensureBenchDeps(b)
	runUnary(b, sctx.UnaryServerInterceptor(), context.Background(), benchInfo(benchMethodAnonymous), &emptypb.Empty{})
}

func BenchmarkFullStack_Anonymous(b *testing.B) {
	ensureBenchDeps(b)
	chain := chainUnaryServer(defaultFrameworkInterceptors(benchAuth)...)
	runUnary(b, chain, fullMetadataCtx(""), benchInfo(benchMethodAnonymous), &emptypb.Empty{})
}

func BenchmarkFullStack_JWT(b *testing.B) {
	ensureBenchDeps(b)
	chain := chainUnaryServer(defaultFrameworkInterceptors(benchAuth)...)
	runUnary(b, chain, fullMetadataCtx(benchJWT), benchInfo(benchMethodAuthed), &emptypb.Empty{})
}

// --- in-process gRPC (bufconn) ---

type benchEchoServer struct{}

func (benchEchoServer) Echo(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

var benchServiceDesc = grpc.ServiceDesc{
	ServiceName: "bench.Service",
	HandlerType: (*interface{})(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Anonymous",
			Handler:    benchUnaryHandler,
		},
		{
			MethodName: "Authed",
			Handler:    benchUnaryHandler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "bench",
}

func benchUnaryHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return (&benchEchoServer{}).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: grpc.ServerTransportStreamFromContext(ctx).Method(),
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return (&benchEchoServer{}).Echo(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func startBenchServer(b *testing.B, unary grpc.UnaryServerInterceptor) (*grpc.ClientConn, func()) {
	b.Helper()
	lis := bufconn.Listen(benchBufSize)
	var opts []grpc.ServerOption
	if unary != nil {
		opts = append(opts, grpc.UnaryInterceptor(unary))
	}
	s := grpc.NewServer(opts...)
	s.RegisterService(&benchServiceDesc, benchEchoServer{})
	go func() {
		_ = s.Serve(lis)
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
	conn, err := grpc.NewClient("passthrough:///bench",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		s.Stop()
		b.Fatal(err)
	}
	cleanup := func() {
		_ = conn.Close()
		s.Stop()
		_ = lis.Close()
	}
	return conn, cleanup
}

func invokeBench(b *testing.B, conn *grpc.ClientConn, method string, md metadata.MD) {
	b.Helper()
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	out := new(emptypb.Empty)
	err := conn.Invoke(ctx, method, &emptypb.Empty{}, out)
	if err != nil {
		b.Fatal(err)
	}
}

func BenchmarkGRPC_NoInterceptor(b *testing.B) {
	ensureBenchDeps(b)
	conn, cleanup := startBenchServer(b, nil)
	defer cleanup()

	md := metadata.Pairs(
		mdpkg.KeyPlatform, fmt.Sprintf("%d", int(common.Platform_PLATFORM_SERVER)),
		mdpkg.KeyLocale, "en",
	)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		invokeBench(b, conn, benchMethodAnonymous, md)
	}
}

func BenchmarkGRPC_FullStack_Anonymous(b *testing.B) {
	ensureBenchDeps(b)
	chain := chainUnaryServer(defaultFrameworkInterceptors(benchAuth)...)
	conn, cleanup := startBenchServer(b, chain)
	defer cleanup()

	pairs := []string{
		mdpkg.KeyProduct, "bench",
		mdpkg.KeyPackage, "pkg",
		mdpkg.KeyVersion, "1",
		mdpkg.KeyPlatform, fmt.Sprintf("%d", int(common.Platform_PLATFORM_SERVER)),
		mdpkg.KeyLocale, "en",
		mdpkg.KeyUserAgent, "bench",
	}
	md := metadata.Pairs(pairs...)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		invokeBench(b, conn, benchMethodAnonymous, md)
	}
}

func BenchmarkGRPC_FullStack_JWT(b *testing.B) {
	ensureBenchDeps(b)
	chain := chainUnaryServer(defaultFrameworkInterceptors(benchAuth)...)
	conn, cleanup := startBenchServer(b, chain)
	defer cleanup()

	md := metadata.Pairs(
		mdpkg.KeyProduct, "bench",
		mdpkg.KeyPackage, "pkg",
		mdpkg.KeyVersion, "1",
		mdpkg.KeyPlatform, fmt.Sprintf("%d", int(common.Platform_PLATFORM_SERVER)),
		mdpkg.KeyLocale, "en",
		mdpkg.KeyUserAgent, "bench",
		mdpkg.KeyToken, benchJWT,
	)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		invokeBench(b, conn, benchMethodAuthed, md)
	}
}
