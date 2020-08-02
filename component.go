package substratum

import (
	"github.com/appootb/protobuf/go/service"
	"github.com/appootb/substratum/discovery"
	"github.com/appootb/substratum/queue"
	"github.com/appootb/substratum/storage"
	"github.com/appootb/substratum/task"
)

// Service component.
type Component interface {
	// Return the component name.
	Name() string

	// Init component.
	Init(discovery.Config) error

	// Init storage.
	InitStorage(storage.Storage) error

	// Init queue consume workers.
	InitQueueWorker(queue.Queue) error

	// Init cron tasks.
	InitCronTask(task.Task) error

	// Register service.
	RegisterService(service.Authenticator, service.Implementor) error
}
