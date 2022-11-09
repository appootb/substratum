package context

import (
	"context"

	"github.com/appootb/substratum/v2/client"
	"github.com/appootb/substratum/v2/discovery"
	ictx "github.com/appootb/substratum/v2/internal/context"
	"github.com/appootb/substratum/v2/logger"
	"github.com/appootb/substratum/v2/queue"
	"github.com/appootb/substratum/v2/service"
	"github.com/appootb/substratum/v2/storage"
	"github.com/appootb/substratum/v2/task"
)

func Context() context.Context {
	return ictx.Context
}

func Cancel() {
	ictx.Cancel()
}

func ServerContext(component string) context.Context {
	return WithServerContext(ictx.Context, component)
}

func WithServerContext(ctx context.Context, component string) context.Context {
	return client.ContextWithConnPool(
		discovery.ContextWithDiscovery(
			logger.ContextWithLogger(
				queue.ContextWithQueueService(
					storage.ContextWithStorage(
						task.ContextWithTaskService(
							service.ContextWithComponentName(ctx, component)))))))
}
