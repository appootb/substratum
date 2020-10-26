package snowflake

import (
	"time"
)

type Sequence interface {
	Next(partition int16, epoch time.Time) (int64, int64, error)
}

type Option func(*Snowflake)

func WithEpoch(epoch time.Time) Option {
	return func(snowflake *Snowflake) {
		snowflake.epoch = epoch
	}
}

func WithPartitionID(partition int64) Option {
	return func(snowflake *Snowflake) {
		snowflake.partition = int16(partition % PartitionIDBitMask)
	}
}

func WithCustomSequence(sequence Sequence) Option {
	return func(snowflake *Snowflake) {
		snowflake.sequence = sequence
	}
}
