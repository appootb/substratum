package credential

import "time"

var (
	clientImpl Client
	serverImpl Server
)

// Return the service implementor.
func ClientImplementor() Client {
	return clientImpl
}

// Register service implementor.
func RegisterClientImplementor(cli Client) {
	clientImpl = cli
}

// Return the service implementor.
func ServerImplementor() Server {
	return serverImpl
}

// Register service implementor.
func RegisterServerImplementor(svr Server) {
	serverImpl = svr
}

// Client secret key.
type Client interface {
	// Add a new secret key.
	Add(accountID uint64, keyID int64, val []byte, expire time.Duration) error
	// Get and refresh the secret key's expiration.
	Refresh(accountID uint64, keyID int64, expire time.Duration) ([]byte, error)
	// Get secret key.
	Get(accountID uint64, keyID int64) ([]byte, error)
	// Revoke the secret key of the specified ID.
	Revoke(accountID uint64, keyID int64) error
	// Revoke all secret keys of the specified account ID.
	RevokeAll(accountID uint64) error
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
