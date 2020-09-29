package substratum

import (
	"net/http"

	"github.com/appootb/protobuf/go/service"
	"github.com/appootb/substratum/configure"
	"github.com/appootb/substratum/queue"
	"github.com/appootb/substratum/storage"
	"github.com/appootb/substratum/task"
)

// Service component.
type Component interface {
	//
	// Invoked when registering component.
	//
	// Return the component name.
	Name() string

	// Init component.
	Init(configure.Configure) error

	// Init storage.
	InitStorage(storage.Storage) error

	// Register HTTP handler.
	RegisterHandler(outer, inner *http.ServeMux) error

	// Register service.
	RegisterService(service.Authenticator, service.Implementor) error

	//
	// Invoked when serving.
	//
	// Run queue consume workers.
	RunQueueWorker(queue.Queue) error

	// Schedule cron tasks.
	ScheduleCronTask(task.Task) error
}
