package model

import (
	"context"

	"github.com/appootb/protobuf/go/common"
	"github.com/appootb/protobuf/go/secret"
	"github.com/appootb/protobuf/go/service"
	"github.com/appootb/substratum/client"
	"github.com/appootb/substratum/logger"
	"github.com/appootb/substratum/metadata"
	"github.com/appootb/substratum/queue"
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

func (m Base) Context() context.Context {
	return m.ctx
}

func (m Base) Metadata() *common.Metadata {
	return metadata.RequestMetadata(m.ctx)
}

func (m Base) AccountSecret() *secret.Info {
	return service.AccountSecretFromContext(m.ctx)
}

func (m Base) Logger() *logger.Helper {
	return logger.ContextLogger(m.ctx)
}

func (m Base) Storage(component string) storage.Storage {
	return storage.ContextStorage(m.ctx, component)
}

func (m Base) ClientConn(component string) *grpc.ClientConn {
	return client.ContextConnPool(m.ctx).Get(component)
}

func (m Base) MessageQueue() queue.Queue {
	return queue.ContextQueueService(m.ctx)
}

func (m Base) CronTask() task.Task {
	return task.ContextTaskService(m.ctx)
}
