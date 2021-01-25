package context

import (
	"context"

	"github.com/appootb/protobuf/go/service"
	"github.com/appootb/substratum/client"
	"github.com/appootb/substratum/discovery"
	"github.com/appootb/substratum/logger"
	"github.com/appootb/substratum/metadata"
	"github.com/appootb/substratum/queue"
	"github.com/appootb/substratum/storage"
	"github.com/appootb/substratum/task"
)

func WithImplementContext(ctx context.Context, component, product string) context.Context {
	return client.ContextWithConnPool(discovery.ContextWithDiscovery(logger.ContextWithLogger(
		queue.ContextWithQueueService(storage.ContextWithStorage(task.ContextWithTaskService(
			metadata.ContextWithProduct(service.ContextWithComponentName(ctx, component), product)))))))
}
