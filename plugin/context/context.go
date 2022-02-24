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

func WithImplementContext(ctx context.Context, component string, product ...string) context.Context {
	if len(product) > 0 && product[0] != "" {
		ctx = metadata.ContextWithProduct(ctx, product[0])
	}
	return client.ContextWithConnPool(
		discovery.ContextWithDiscovery(
			logger.ContextWithLogger(
				queue.ContextWithQueueService(
					storage.ContextWithStorage(
						task.ContextWithTaskService(
							service.ContextWithComponentName(ctx, component)))))))
}
