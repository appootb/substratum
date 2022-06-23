package context

import (
	"context"

	"github.com/appootb/substratum/v2/client"
	"github.com/appootb/substratum/v2/discovery"
	"github.com/appootb/substratum/v2/logger"
	"github.com/appootb/substratum/v2/queue"
	"github.com/appootb/substratum/v2/service"
	"github.com/appootb/substratum/v2/storage"
	"github.com/appootb/substratum/v2/task"
)

func WithImplementContext(ctx context.Context, component string) context.Context {
	return client.ContextWithConnPool(
		discovery.ContextWithDiscovery(
			logger.ContextWithLogger(
				queue.ContextWithQueueService(
					storage.ContextWithStorage(
						task.ContextWithTaskService(
							service.ContextWithComponentName(ctx, component)))))))
}
