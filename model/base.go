package model

import (
	"context"

	"github.com/appootb/substratum/v2/client"
	"github.com/appootb/substratum/v2/discovery"
	"github.com/appootb/substratum/v2/logger"
	"github.com/appootb/substratum/v2/metadata"
	"github.com/appootb/substratum/v2/proto/go/common"
	"github.com/appootb/substratum/v2/proto/go/secret"
	"github.com/appootb/substratum/v2/queue"
	"github.com/appootb/substratum/v2/service"
	"github.com/appootb/substratum/v2/storage"
	"github.com/appootb/substratum/v2/task"
	"google.golang.org/grpc"
)

type Base struct {
	ctx context.Context
}

func New(opts ...Option) Base {
	base := Base{}
	for _, opt := range opts {
		opt(&base)
	}
	return base
}

func (m Base) RpcContext(keyID int64) context.Context {
	return client.WithContext(m.Context(), keyID)
}

func (m Base) Context() context.Context {
	return m.ctx
}

func (m Base) Metadata() *common.Metadata {
	return metadata.IncomingMetadata(m.ctx)
}

func (m Base) AccountSecret() *secret.Info {
	return service.AccountSecretFromContext(m.ctx)
}

func (m Base) Discovery() discovery.Discovery {
	return discovery.ContextDiscovery(m.ctx)
}

func (m Base) Logger() *logger.Helper {
	return logger.ContextLogger(m.ctx)
}

func (m Base) Storage() storage.Storage {
	component := service.ComponentNameFromContext(m.Context())
	return storage.ContextStorage(m.ctx, component)
}

func (m Base) ClientConn(target string) *grpc.ClientConn {
	return client.ContextConnPool(m.ctx).Get(target)
}

func (m Base) MessageQueue() queue.Queue {
	return queue.ContextQueueService(m.ctx)
}

func (m Base) CronTask() task.Task {
	return task.ContextTaskService(m.ctx)
}
