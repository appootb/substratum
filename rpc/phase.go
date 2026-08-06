package rpc

import (
	"fmt"

	"github.com/appootb/substratum/v2/auth"
	"github.com/appootb/substratum/v2/client"
	"github.com/appootb/substratum/v2/discovery"
	"github.com/appootb/substratum/v2/errors"
	"github.com/appootb/substratum/v2/logger"
	"github.com/appootb/substratum/v2/metadata"
	"github.com/appootb/substratum/v2/monitor"
	"github.com/appootb/substratum/v2/queue"
	"github.com/appootb/substratum/v2/recovery"
	"github.com/appootb/substratum/v2/storage"
	"github.com/appootb/substratum/v2/task"
	"google.golang.org/grpc"

	validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
)

// UnaryPhase selects where custom unary interceptors are merged into the
// default server chain. Register from plugin Init or before substratum.NewServer.
type UnaryPhase int

const (
	UnaryOutermost UnaryPhase = iota
	UnaryBeforeAuth
	UnaryAfterAuth
	UnaryBeforeHandler
)

// StreamPhase selects where custom stream interceptors are merged.
type StreamPhase int

const (
	StreamOutermost StreamPhase = iota
	StreamBeforeAuth
	StreamAfterAuth
	StreamBeforeHandler
)

const unaryPhaseCount = int(UnaryBeforeHandler) + 1
const streamPhaseCount = int(StreamBeforeHandler) + 1

var (
	unaryPhases  [unaryPhaseCount][]grpc.UnaryServerInterceptor
	streamPhases [streamPhaseCount][]grpc.StreamServerInterceptor

	frozenUnaryChain  []grpc.UnaryServerInterceptor
	frozenStreamChain []grpc.StreamServerInterceptor
)

// RegisterUnary appends unary interceptors to a phase slot.
func RegisterUnary(phase UnaryPhase, fns ...grpc.UnaryServerInterceptor) {
	if len(fns) == 0 {
		return
	}
	registerUnary(phase, fns, false)
}

// SetUnary replaces all interceptors registered for a unary phase.
func SetUnary(phase UnaryPhase, fns ...grpc.UnaryServerInterceptor) {
	registerUnary(phase, fns, true)
}

// RegisterStream appends stream interceptors to a phase slot.
func RegisterStream(phase StreamPhase, fns ...grpc.StreamServerInterceptor) {
	if len(fns) == 0 {
		return
	}
	registerStream(phase, fns, false)
}

// SetStream replaces all interceptors registered for a stream phase.
func SetStream(phase StreamPhase, fns ...grpc.StreamServerInterceptor) {
	registerStream(phase, fns, true)
}

// BuildInterceptorChains merges phase interceptors with framework defaults.
// Called from plugin.Register; freezes further registration.
func BuildInterceptorChains() {
	BuildUnaryChain()
	BuildStreamChain()
}

// UnaryChain returns the built unary interceptor chain.
func UnaryChain() []grpc.UnaryServerInterceptor {
	if frozenUnaryChain == nil {
		panic("rpc: unary interceptor chain not built")
	}
	return frozenUnaryChain
}

// StreamChain returns the built stream interceptor chain.
func StreamChain() []grpc.StreamServerInterceptor {
	if frozenStreamChain == nil {
		panic("rpc: stream interceptor chain not built")
	}
	return frozenStreamChain
}

// BuildUnaryChain builds the merged unary chain once and returns it.
func BuildUnaryChain() []grpc.UnaryServerInterceptor {
	if frozenUnaryChain != nil {
		return frozenUnaryChain
	}
	chain := make([]grpc.UnaryServerInterceptor, 0, 32)
	appendPhase := func(phase UnaryPhase) {
		chain = append(chain, unaryPhases[phase]...)
	}
	appendPhase(UnaryOutermost)
	chain = append(chain,
		recovery.UnaryServerInterceptor(),
		metadata.UnaryServerInterceptor(),
		errors.UnaryResponseInterceptor(),
		monitor.UnaryServerInterceptor(),
	)
	appendPhase(UnaryBeforeAuth)
	chain = append(chain, auth.UnaryServerInterceptor())
	appendPhase(UnaryAfterAuth)
	chain = append(chain,
		logger.UnaryServerInterceptor(),
		validator.UnaryServerInterceptor(),
		client.UnaryServerInterceptor(),
		discovery.UnaryServerInterceptor(),
		storage.UnaryServerInterceptor(),
		queue.UnaryServerInterceptor(),
		task.UnaryServerInterceptor(),
	)
	appendPhase(UnaryBeforeHandler)
	frozenUnaryChain = chain
	return frozenUnaryChain
}

// BuildStreamChain builds the merged stream chain once and returns it.
func BuildStreamChain() []grpc.StreamServerInterceptor {
	if frozenStreamChain != nil {
		return frozenStreamChain
	}
	chain := make([]grpc.StreamServerInterceptor, 0, 32)
	appendPhase := func(phase StreamPhase) {
		chain = append(chain, streamPhases[phase]...)
	}
	appendPhase(StreamOutermost)
	chain = append(chain,
		recovery.StreamServerInterceptor(),
		metadata.StreamServerInterceptor(),
		errors.StreamServerInterceptor(),
		monitor.StreamServerInterceptor(),
	)
	appendPhase(StreamBeforeAuth)
	chain = append(chain, auth.StreamServerInterceptor())
	appendPhase(StreamAfterAuth)
	chain = append(chain,
		logger.StreamServerInterceptor(),
		validator.StreamServerInterceptor(),
		client.StreamServerInterceptor(),
		discovery.StreamServerInterceptor(),
		storage.StreamServerInterceptor(),
		queue.StreamServerInterceptor(),
		task.StreamServerInterceptor(),
	)
	appendPhase(StreamBeforeHandler)
	frozenStreamChain = chain
	return frozenStreamChain
}

func registerUnary(phase UnaryPhase, fns []grpc.UnaryServerInterceptor, replace bool) {
	if frozenUnaryChain != nil {
		panic(fmt.Sprintf("rpc: RegisterUnary(%v) after interceptor chain built", phase))
	}
	if replace {
		unaryPhases[phase] = append([]grpc.UnaryServerInterceptor(nil), fns...)
		return
	}
	unaryPhases[phase] = append(unaryPhases[phase], fns...)
}

func registerStream(phase StreamPhase, fns []grpc.StreamServerInterceptor, replace bool) {
	if frozenStreamChain != nil {
		panic(fmt.Sprintf("rpc: RegisterStream(%v) after interceptor chain built", phase))
	}
	if replace {
		streamPhases[phase] = append([]grpc.StreamServerInterceptor(nil), fns...)
		return
	}
	streamPhases[phase] = append(streamPhases[phase], fns...)
}

// resetPhases clears phase registration and built chains. For tests only.
func resetPhases() {
	unaryPhases = [unaryPhaseCount][]grpc.UnaryServerInterceptor{}
	streamPhases = [streamPhaseCount][]grpc.StreamServerInterceptor{}
	frozenUnaryChain = nil
	frozenStreamChain = nil
}
