package credential

import "time"

var (
	clientImpl Client
	serverImpl Server
)

// ClientImplementor returns the client-side seed service implementor.
func ClientImplementor() Client {
	return clientImpl
}

// RegisterClientImplementor registers the client-side seed service implementor.
func RegisterClientImplementor(cli Client) {
	clientImpl = cli
}

// ServerImplementor returns the server-side seed service implementor.
func ServerImplementor() Server {
	return serverImpl
}

// RegisterServerImplementor registers the server-side seed service implementor.
func RegisterServerImplementor(svr Server) {
	serverImpl = svr
}

// Client secret key.
type Client interface {
	// Add creates a new secret key.
	Add(uid uint64, keyID int64, val []byte, expire time.Duration) error
	// Refresh gets and refreshes the secret key's expiration.
	Refresh(uid uint64, keyID int64, expire time.Duration) ([]byte, error)
	// Get secret key.
	Get(uid uint64, keyID int64) ([]byte, error)
	// Revoke removes the secret key of the specified ID.
	Revoke(uid uint64, keyID int64) error
	// RevokeAll removes all secret keys of the specified user ID.
	RevokeAll(uid uint64) error
	// Lock disables all secret keys for a specified duration.
	// Returns codes.FailedPrecondition (9).
	Lock(uid uint64, reason string, duration time.Duration) error
	// Unlock enables the secret keys.
	Unlock(uid uint64) error
}

// Server secret key.
type Server interface {
	// Add a new secret key for the specified ID.
	Add(keyID int64, key []byte) error
	// Get the secret key of the specified ID.
	Get(keyID int64) ([]byte, error)
	// Revoke the secret key of the specified ID.
	Revoke(keyID int64) error
}
