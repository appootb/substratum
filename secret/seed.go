package secret

import (
	"fmt"
)

var DefaultSeed = newMemorySeed()

// Seed
type Seed interface {
	// Generate a new seed.
	New(masterId, slaveId interface{}) ([]byte, error)
	// Get seed.
	Get(masterId, slaveId interface{}) ([]byte, error)
	// Revoke the seed of the specified ID.
	Revoke(masterId, slaveId interface{}) error
	// Revoke all seed of the specified master ID.
	RevokeAll(masterId interface{}) error
}

type memory struct{}

// Used only for debug.
func newMemorySeed() Seed {
	return &memory{}
}

func (s *memory) New(masterId, slaveId interface{}) ([]byte, error) {
	return []byte(fmt.Sprintf("%v-%v", masterId, slaveId)), nil
}

func (s *memory) Get(masterId, slaveId interface{}) ([]byte, error) {
	return []byte(fmt.Sprintf("%v-%v", masterId, slaveId)), nil
}

func (s *memory) Revoke(_, _ interface{}) error {
	return nil
}

func (s *memory) RevokeAll(_ interface{}) error {
	return nil
}
