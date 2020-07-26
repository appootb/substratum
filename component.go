package substratum

import (
	"github.com/appootb/protobuf/go/service"
	"github.com/appootb/substratum/cron"
	"github.com/appootb/substratum/discovery"
	"github.com/appootb/substratum/queue"
	"github.com/appootb/substratum/storage"
)

// Service component.
type Component interface {
	// Return the component name.
	Name() string

	// Init component.
	Init(discovery.Config) error

	// Init storage.
	InitStorage(storage.Storage) error

	// Init queue job consume workers.
	InitJobWorker(queue.Service) error

	// Init cron tasks.
	InitCron(cron.Cron) error

	// Register service.
	RegisterService(service.Authenticator, service.Implementor) error
}
