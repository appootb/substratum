package snowflake

import "time"

type Option func(*Snowflake)

func WithEpoch(epoch time.Time) Option {
	return func(snowflake *Snowflake) {
		snowflake.epoch = epoch
	}
}

func WithNodeID(node int16) Option {
	return func(snowflake *Snowflake) {
		snowflake.node = node
	}
}
