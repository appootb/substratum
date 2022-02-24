package model

import (
	"context"

	"github.com/appootb/substratum/client"
	"github.com/appootb/substratum/discovery"
	"github.com/appootb/substratum/logger"
	"github.com/appootb/substratum/metadata"
	"github.com/appootb/substratum/proto/go/common"
	"github.com/appootb/substratum/proto/go/secret"
	"github.com/appootb/substratum/queue"
	"github.com/appootb/substratum/service"
	"github.com/appootb/substratum/storage"
	"github.com/appootb/substratum/task"
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

func (m Base) RpcContext(keyID int64, product ...string) context.Context {
	return client.WithContext(m.Context(), keyID, product...)
}

func (m Base) Context() context.Context {
	return m.ctx
}

func (m Base) Metadata() *common.Metadata {
	return metadata.RequestMetadata(m.ctx)
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
