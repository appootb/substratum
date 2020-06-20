package model

import (
	"context"

	"github.com/appootb/protobuf/go/common"
	"github.com/appootb/protobuf/go/permission"
	"github.com/appootb/protobuf/go/service"
	"github.com/appootb/substratum/logger"
	"github.com/appootb/substratum/metadata"
	"github.com/appootb/substratum/storage"
)

type Base struct {
	ctx context.Context
}

func New(opts ...Option) Base {
	base := Base{}
	for _, opt := range opts {
		opt(&base)
	}
	return base
}

func (m Base) Context() context.Context {
	return m.ctx
}

func (m Base) Metadata() *common.Metadata {
	return metadata.RequestMetadata(m.ctx)
}

func (m Base) AccountSecret() *permission.Secret {
	return service.AccountSecretFromContext(m.ctx)
}

func (m Base) Logger() *logger.Helper {
	return logger.ContextLogger(m.ctx)
}

func (m Base) Storage(component string) storage.Storage {
	return storage.ContextStorage(m.ctx, component)
}
