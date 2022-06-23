package plugin

import (
	"sync"

	"github.com/appootb/substratum/v2/plugin/auth"
	"github.com/appootb/substratum/v2/plugin/client"
	"github.com/appootb/substratum/v2/plugin/configure"
	"github.com/appootb/substratum/v2/plugin/credential"
	"github.com/appootb/substratum/v2/plugin/discovery"
	"github.com/appootb/substratum/v2/plugin/errors"
	"github.com/appootb/substratum/v2/plugin/logger"
	"github.com/appootb/substratum/v2/plugin/queue"
	"github.com/appootb/substratum/v2/plugin/resolver"
	"github.com/appootb/substratum/v2/plugin/storage"
	"github.com/appootb/substratum/v2/plugin/task"
	"github.com/appootb/substratum/v2/plugin/token"
)

var (
	once sync.Once
)

func Register() {
	once.Do(func() {
		// Client ConnPool
		client.Init()
		// Config
		configure.Init()
		// Credential
		credential.Init()
		// Discovery
		discovery.Init()
		// Errors
		errors.Init()
		// Logger
		logger.Init()
		// Queue
		queue.Init()
		// Resolver
		resolver.Init()
		// Storage
		storage.Init()
		// Task
		task.Init()
		// Token
		token.Init()
		// Auth
		auth.Init()
	})
}
