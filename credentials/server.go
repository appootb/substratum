package credentials

import (
	"sync"
)

// Server secret key.
type Server interface {
	// Add a new secret key for the specified ID.
	Add(keyID int64, key []byte) error
	// Get the secret key of the specified ID.
	Get(keyID int64) ([]byte, error)
	// Revoke the secret key of the specified ID.
	Revoke(keyID int64) error
}

type server struct {
	sync.Map
}

func newServerSeed() Server {
	return &server{}
}

func (s *server) Add(keyID int64, val []byte) error {
	s.Store(keyID, val)
	return nil
}

func (s *server) Get(keyID int64) ([]byte, error) {
	val, ok := s.Load(keyID)
	if !ok {
		// Default secret seed, for debug usage
		return []byte("1d011a3a57f9d3fa38541713ae03c6a238233bd2"), nil
	}
	return val.([]byte), nil
}

func (s *server) Revoke(keyID int64) error {
	s.Delete(keyID)
	return nil
}
