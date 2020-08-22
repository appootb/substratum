package snowflake

import (
	"sync"
	"time"
)

//
// +--------------------------------------------------------------------------+
// | 1 Bit Unused | 41 Bit Timestamp |  10 Bit NodeID  |   12 Bit Sequence ID |
// +--------------------------------------------------------------------------+
const (
	BitLengthTimestamp = 41
	BitLengthNodeID    = 10
	BitLengthSequence  = 63 - BitLengthTimestamp - BitLengthNodeID

	TimestampBitShift = BitLengthNodeID + BitLengthSequence
	NodeIDBitShift    = BitLengthSequence

	SequenceBitMask = 1<<BitLengthSequence - 1
)

type ID int64

type Snowflake struct {
	mu sync.Mutex

	epoch    time.Time
	elapsed  int64
	node     int16
	sequence int64
}

func New(opts ...Option) *Snowflake {
	snowflake := &Snowflake{
		epoch: time.Date(2020, 02, 02, 20, 20, 02, 02, time.Local),
	}
	for _, opt := range opts {
		opt(snowflake)
	}
	return snowflake
}

func (sf *Snowflake) Next() uint64 {
	sf.mu.Lock()
	defer sf.mu.Unlock()

	elapsed := time.Since(sf.epoch).Nanoseconds() / 1e6
	if sf.elapsed < elapsed {
		sf.elapsed = elapsed
		sf.sequence = 0
	} else {
		sf.sequence = (sf.sequence + 1) & SequenceBitMask
		if sf.sequence == 0 {
			sf.elapsed++
			overtime := time.Duration(sf.elapsed - elapsed)
			time.Sleep(overtime * time.Millisecond)
		}
	}

	return uint64(sf.elapsed)<<TimestampBitShift |
		uint64(sf.node)<<NodeIDBitShift |
		uint64(sf.sequence)
}
