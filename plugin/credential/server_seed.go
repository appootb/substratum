package credential

import "sync"

type ServerSeed struct {
	sync.Map
}

func (s *ServerSeed) Add(keyID int64, val []byte) error {
	s.Store(keyID, val)
	return nil
}

func (s *ServerSeed) Get(keyID int64) ([]byte, error) {
	val, ok := s.Load(keyID)
	if !ok {
		// Default secret seed, for debug usage
		return []byte("1d011a3a57f9d3fa38541713ae03c6a238233bd2"), nil
	}
	return val.([]byte), nil
}

func (s *ServerSeed) Revoke(keyID int64) error {
	s.Delete(keyID)
	return nil
}
