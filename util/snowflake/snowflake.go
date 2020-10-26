package snowflake

import (
	"sync"
	"time"
)

//
// +-------------------------------------------------------------------------------+
// | 1 Bit Unused | 41 Bit Timestamp |  10 Bit PartitionID  |   12 Bit Sequence ID |
// +-------------------------------------------------------------------------------+
const (
	BitLengthTimestamp   = 41
	BitLengthPartitionID = 10
	BitLengthSequence    = 63 - BitLengthTimestamp - BitLengthPartitionID

	TimestampBitShift   = BitLengthPartitionID + BitLengthSequence
	PartitionIDBitShift = BitLengthSequence

	PartitionIDBitMask = 1<<BitLengthPartitionID - 1
	SequenceBitMask    = 1<<BitLengthSequence - 1
)

var Default = New()

func SetPartitionID(partitionID int64) {
	Default.mu.Lock()
	Default.partition = int16(partitionID % PartitionIDBitMask)
	Default.mu.Unlock()
}

func NextID() uint64 {
	return Default.Next()
}

type Snowflake struct {
	mu sync.RWMutex

	epoch     time.Time
	partition int16
	sequence  Sequence
}

func New(opts ...Option) *Snowflake {
	snowflake := &Snowflake{
		epoch:    time.Date(2020, 02, 02, 20, 20, 02, 02, time.Local),
		sequence: &DefaultSequence{},
	}
	for _, opt := range opts {
		opt(snowflake)
	}
	return snowflake
}

func (sf *Snowflake) Next() uint64 {
	sf.mu.RLock()
	defer sf.mu.RUnlock()

	return sf.sequence.Next(sf.partition, sf.epoch)
}

func (sf *Snowflake) Timestamp(id uint64) time.Time {
	dur := id >> (BitLengthPartitionID + BitLengthSequence)
	return sf.epoch.Add(time.Duration(dur) * time.Millisecond)
}
