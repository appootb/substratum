package plugin

import (
	"sync"

	"github.com/appootb/substratum/plugin/auth"
	"github.com/appootb/substratum/plugin/balancer"
	"github.com/appootb/substratum/plugin/client"
	"github.com/appootb/substratum/plugin/configure"
	"github.com/appootb/substratum/plugin/credential"
	"github.com/appootb/substratum/plugin/discovery"
	"github.com/appootb/substratum/plugin/errors"
	"github.com/appootb/substratum/plugin/logger"
	"github.com/appootb/substratum/plugin/queue"
	"github.com/appootb/substratum/plugin/resolver"
	"github.com/appootb/substratum/plugin/storage"
	"github.com/appootb/substratum/plugin/task"
	"github.com/appootb/substratum/plugin/token"
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
		// Balancer
		balancer.Init()
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
