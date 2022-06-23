package substratum

import (
	"net/http"

	"github.com/appootb/substratum/v2/configure"
	"github.com/appootb/substratum/v2/queue"
	"github.com/appootb/substratum/v2/service"
	"github.com/appootb/substratum/v2/storage"
	"github.com/appootb/substratum/v2/task"
)

// Component interface.
type Component interface {
	//
	// Invoked when registering component.
	//

	// Name returns the component name.
	Name() string

	// Init component.
	Init(configure.Configure) error

	// InitStorage entry.
	InitStorage(storage.Storage) error

	// RegisterHandler registers HTTP handler.
	RegisterHandler(outer, inner http.Handler) error

	// RegisterService registers the service.
	RegisterService(service.Authenticator, service.Implementor) error

	//
	// Invoked when serving.
	//

	// RunQueueWorker will run the registered queue consumers.
	RunQueueWorker(queue.Queue) error

	// ScheduleCronTask schedules the cron tasks.
	ScheduleCronTask(task.Task) error
}
