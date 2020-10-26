package snowflake

import (
	"sync"
	"time"
)

type DefaultSequence struct {
	mu sync.Mutex

	elapsed  int64
	sequence int64
}

func (s *DefaultSequence) Next(partition int16, epoch time.Time) uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	elapsed := time.Since(epoch).Nanoseconds() / 1e6
	if s.elapsed < elapsed {
		s.elapsed = elapsed
		s.sequence = 0
	} else {
		s.sequence = (s.sequence + 1) & SequenceBitMask
		if s.sequence == 0 {
			s.elapsed++
			overtime := time.Duration(s.elapsed - elapsed)
			time.Sleep(overtime * time.Millisecond)
		}
	}

	return uint64(s.elapsed)<<TimestampBitShift |
		uint64(partition)<<PartitionIDBitShift |
		uint64(s.sequence)
}
