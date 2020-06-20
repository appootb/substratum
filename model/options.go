package model

import "context"

type Option func(*Base)

func WithContext(ctx context.Context) Option {
	return func(base *Base) {
		base.ctx = ctx
	}
}
