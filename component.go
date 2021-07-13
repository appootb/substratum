package substratum

import (
	"net/http"

	"github.com/appootb/substratum/configure"
	"github.com/appootb/substratum/queue"
	"github.com/appootb/substratum/service"
	"github.com/appootb/substratum/storage"
	"github.com/appootb/substratum/task"
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
	RegisterHandler(outer, inner *http.ServeMux) error

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
