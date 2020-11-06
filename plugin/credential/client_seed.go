package credential

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

type ClientSeed struct {
	sync.Map
}

func (s *ClientSeed) Add(accountID uint64, keyID int64, val []byte, _ time.Duration) error {
	key := fmt.Sprintf("%d-%d", accountID, keyID)
	s.Store(key, val)
	return nil
}

func (s *ClientSeed) Refresh(accountID uint64, keyID int64, _ time.Duration) ([]byte, error) {
	return s.Get(accountID, keyID)
}

func (s *ClientSeed) Get(accountID uint64, keyID int64) ([]byte, error) {
	key := fmt.Sprintf("%d-%d", accountID, keyID)
	val, ok := s.Load(key)
	if !ok {
		return nil, errors.New("substratum: client key not found:" + key)
	}
	return val.([]byte), nil
}

func (s *ClientSeed) Revoke(accountID uint64, keyID int64) error {
	key := fmt.Sprintf("%d-%d", accountID, keyID)
	s.Delete(key)
	return nil
}

func (s *ClientSeed) RevokeAll(accountID uint64) error {
	key := fmt.Sprintf("%d-", accountID)
	s.Range(func(k, _ interface{}) bool {
		if strings.HasPrefix(k.(string), key) {
			s.Delete(k)
		}
		return true
	})
	return nil
}
