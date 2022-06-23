package credential

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/appootb/substratum/v2/errors"
	"google.golang.org/grpc/codes"
)

type clientSeedInfo struct {
	PrivateKey  []byte
	NotBefore   time.Time
	NotAfter    time.Time
	LockMessage string
}

type ClientSeed struct {
	sync.Map
}

func (s *ClientSeed) Add(accountID uint64, keyID int64, val []byte, expire time.Duration) error {
	key := fmt.Sprintf("%d-%d", accountID, keyID)
	s.Store(key, &clientSeedInfo{
		PrivateKey: val,
		NotAfter:   time.Now().Add(expire),
	})
	return nil
}

func (s *ClientSeed) Refresh(accountID uint64, keyID int64, _ time.Duration) ([]byte, error) {
	return s.Get(accountID, keyID)
}

func (s *ClientSeed) Get(accountID uint64, keyID int64) ([]byte, error) {
	key := fmt.Sprintf("%d-%d", accountID, keyID)
	val, ok := s.Load(key)
	if !ok {
		return nil, errors.New(codes.Unauthenticated, "substratum: client key not found:"+key)
	}
	info := val.(*clientSeedInfo)
	if time.Now().After(info.NotAfter) {
		return nil, errors.New(codes.Unauthenticated, "substratum: client key expired")
	}
	if !info.NotBefore.IsZero() && time.Now().Before(info.NotBefore) {
		return nil, errors.New(codes.FailedPrecondition, info.LockMessage)
	}
	return info.PrivateKey, nil
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

func (s *ClientSeed) Lock(accountID uint64, reason string, duration time.Duration) error {
	key := fmt.Sprintf("%d-", accountID)
	s.Range(func(k, v interface{}) bool {
		if strings.HasPrefix(k.(string), key) {
			info := v.(*clientSeedInfo)
			info.NotBefore = time.Now().Add(duration)
			info.LockMessage = reason
			s.Store(k, info)
		}
		return true
	})
	return nil
}

func (s *ClientSeed) Unlock(accountID uint64) error {
	key := fmt.Sprintf("%d-", accountID)
	s.Range(func(k, v interface{}) bool {
		if strings.HasPrefix(k.(string), key) {
			info := v.(*clientSeedInfo)
			info.NotBefore = time.Date(1, 1, 1, 0, 0, 0, 0, time.Local)
			info.LockMessage = ""
			s.Store(k, info)
		}
		return true
	})
	return nil
}
