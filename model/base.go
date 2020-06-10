package model

import (
	"context"

	"github.com/appootb/protobuf/go/common"
	"github.com/appootb/protobuf/go/permission"
	"github.com/appootb/protobuf/go/service"
	"github.com/appootb/substratum/logger"
	"github.com/appootb/substratum/metadata"
)

type Base struct {
	ctx context.Context
}

func New(ctx context.Context) Base {
	return Base{
		ctx: ctx,
	}
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
