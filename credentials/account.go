package credentials

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

// Client secret key.
type Client interface {
	// Add a new secret key.
	Add(accountID uint64, keyID int64, val []byte) error
	// Get secret key.
	Get(accountID uint64, keyID int64) ([]byte, error)
	// Revoke the secret key of the specified ID.
	Revoke(accountID uint64, keyID int64) error
	// Revoke all secret keys of the specified account ID.
	RevokeAll(accountID uint64) error
}

type client struct {
	sync.Map
}

// Used only for debug.
func newClientSeed() Client {
	return &client{}
}

func (s *client) Add(accountID uint64, keyID int64, val []byte) error {
	key := fmt.Sprintf("%d-%d", accountID, keyID)
	s.Store(key, val)
	return nil
}

func (s *client) Get(accountID uint64, keyID int64) ([]byte, error) {
	key := fmt.Sprintf("%d-%d", accountID, keyID)
	val, ok := s.Load(key)
	if !ok {
		return nil, errors.New("substratum: client key not found:" + key)
	}
	return val.([]byte), nil
}

func (s *client) Revoke(accountID uint64, keyID int64) error {
	key := fmt.Sprintf("%d-%d", accountID, keyID)
	s.Delete(key)
	return nil
}

func (s *client) RevokeAll(accountID uint64) error {
	key := fmt.Sprintf("%d-", accountID)
	s.Range(func(k, _ interface{}) bool {
		if strings.HasPrefix(k.(string), key) {
			s.Delete(k)
		}
		return true
	})
	return nil
}
