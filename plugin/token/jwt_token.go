package token

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/appootb/substratum/credential"
	"github.com/appootb/substratum/proto/go/common"
	"github.com/appootb/substratum/proto/go/permission"
	"github.com/appootb/substratum/proto/go/secret"
	"github.com/appootb/substratum/token"
	"github.com/appootb/substratum/util/datetime"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/gbrlsnchs/jwt/v3/jwtutil"
)

var (
	UnsupportedAlgorithm = errors.New("substratum: algorithm not supported")
)

func Init() {
	if token.Implementor() == nil {
		token.RegisterImplementor(&JwtToken{})
	}
}

type JwtToken struct{}

func (t *JwtToken) sign(s *secret.Info, key []byte) (string, error) {
	issuedAt := datetime.FromProtoTime(s.GetIssuedAt()).Time
	payload := &jwt.Payload{
		Issuer:         s.GetIssuer(),
		Subject:        s.GetSubject().String(),
		Audience:       s.GetRoles(),
		ExpirationTime: jwt.NumericDate(datetime.FromProtoTime(s.GetExpiredAt()).Time),
		NotBefore:      jwt.NumericDate(issuedAt.Add(-time.Minute)),
		IssuedAt:       jwt.NumericDate(issuedAt),
		JWTID:          strconv.FormatUint(s.GetAccount(), 10),
	}
	alg, err := t.getAlgorithm(s.GetAlgorithm(), key)
	if err != nil {
		return "", err
	}
	contentType := jwt.ContentType(s.GetType().String())
	keyID := jwt.KeyID(fmt.Sprintf("%v-%v", s.GetAccount(), s.GetKeyId()))
	val, err := jwt.Sign(payload, alg, contentType, keyID)
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func (t *JwtToken) getAlgorithm(alg secret.Algorithm, key []byte) (jwt.Algorithm, error) {
	switch alg {
	case secret.Algorithm_HMAC:
		return jwt.NewHS512(key), nil
	case secret.Algorithm_RSA:
		priv, err := x509.ParsePKCS1PrivateKey(key)
		if err != nil {
			return nil, err
		}
		return jwt.NewRS512(jwt.RSAPrivateKey(priv), jwt.RSAPublicKey(&priv.PublicKey)), nil
	case secret.Algorithm_PSS:
		priv, err := x509.ParsePKCS1PrivateKey(key)
		if err != nil {
			return nil, err
		}
		return jwt.NewPS512(jwt.RSAPrivateKey(priv), jwt.RSAPublicKey(&priv.PublicKey)), nil
	case secret.Algorithm_ECDSA:
		priv, err := x509.ParseECPrivateKey(key)
		if err != nil {
			return nil, err
		}
		return jwt.NewES512(jwt.ECDSAPrivateKey(priv), jwt.ECDSAPublicKey(&priv.PublicKey)), nil
	case secret.Algorithm_EdDSA:
		priv := ed25519.PrivateKey(key)
		return jwt.NewEd25519(jwt.Ed25519PrivateKey(priv), jwt.Ed25519PublicKey(priv.Public().([]byte))), nil
	default:
		return nil, UnsupportedAlgorithm
	}
}

func (t *JwtToken) NewSecretKey(alg secret.Algorithm) ([]byte, error) {
	switch alg {
	case secret.Algorithm_HMAC:
		var key [16]byte
		_, err := rand.Read(key[:])
		return key[:], err
	case secret.Algorithm_RSA, secret.Algorithm_PSS:
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, err
		}
		return x509.MarshalPKCS1PrivateKey(key), nil
	case secret.Algorithm_ECDSA:
		key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		if err != nil {
			return nil, err
		}
		return x509.MarshalECPrivateKey(key)
	case secret.Algorithm_EdDSA:
		_, key, err := ed25519.GenerateKey(rand.Reader)
		return key, err
	default:
		return nil, UnsupportedAlgorithm
	}
}

// Generate a new token with specified options.
func (t *JwtToken) Generate(s *secret.Info) (string, error) {
	if s.GetType() == secret.Type_SERVER {
		// TODO: Server key should be added by operators.
		key, err := credential.ServerImplementor().Get(s.GetKeyId())
		if err != nil {
			return "", err
		}
		return t.sign(s, key)
	}

	key, err := t.NewSecretKey(s.GetAlgorithm())
	if err != nil {
		return "", err
	}
	dur := datetime.FromProtoTime(s.GetExpiredAt()).Time.Sub(datetime.FromProtoTime(s.GetIssuedAt()).Time)
	if err := credential.ClientImplementor().Add(s.GetAccount(), s.GetKeyId(), key, dur); err != nil {
		return "", err
	}
	return t.sign(s, key)
}

// Refresh the token with expired time renewed.
func (t *JwtToken) Refresh(s *secret.Info) (string, error) {
	var (
		err error
		key []byte
	)
	if s.GetType() == secret.Type_SERVER {
		key, err = credential.ServerImplementor().Get(s.GetKeyId())
	} else {
		dur := datetime.FromProtoTime(s.GetExpiredAt()).Time.Sub(datetime.FromProtoTime(s.GetIssuedAt()).Time)
		key, err = credential.ClientImplementor().Refresh(s.GetAccount(), s.GetKeyId(), dur)
	}
	if err != nil {
		return "", err
	}
	// Calculate new timestamp.
	issuedAt := datetime.FromProtoTime(s.GetIssuedAt()).Time
	expiredAt := datetime.FromProtoTime(s.GetExpiredAt()).Time
	now := time.Now()
	s.IssuedAt = datetime.WithTime(now).Proto()
	s.ExpiredAt = datetime.WithTime(now.Add(expiredAt.Sub(issuedAt))).Proto()
	return t.sign(s, key)
}

// Parse the metadata.
func (t *JwtToken) Parse(md *common.Metadata) (*secret.Info, error) {
	return t.ParseRaw(md.GetToken())
}

// ParseRaw parses a token string.
func (t *JwtToken) ParseRaw(token string) (*secret.Info, error) {
	var (
		accountID uint64
		keyID     int64
		alg       secret.Algorithm
	)
	resolver := &jwtutil.Resolver{
		New: func(header jwt.Header) (jwt.Algorithm, error) {
			var (
				err error
				key []byte
			)
			keyIDs := strings.Split(header.KeyID, "-")
			if len(keyIDs) != 2 {
				return nil, jwt.ErrAlgValidation
			}
			accountID, _ = strconv.ParseUint(keyIDs[0], 10, 64)
			keyID, _ = strconv.ParseInt(keyIDs[1], 10, 64)
			if header.ContentType == secret.Type_SERVER.String() {
				key, err = credential.ServerImplementor().Get(keyID)
			} else {
				key, err = credential.ClientImplementor().Get(accountID, keyID)
			}
			if err != nil {
				return nil, err
			}
			switch header.Algorithm {
			case "HS512":
				alg = secret.Algorithm_HMAC
			case "RS512":
				alg = secret.Algorithm_RSA
			case "PS512":
				alg = secret.Algorithm_PSS
			case "ES512":
				alg = secret.Algorithm_ECDSA
			case "Ed25519":
				alg = secret.Algorithm_EdDSA
			default:
				alg = secret.Algorithm_None
			}
			return t.getAlgorithm(alg, key)
		},
	}
	// Verify
	now := time.Now()
	payload := jwt.Payload{}
	validator := jwt.ValidatePayload(&payload,
		jwt.NotBeforeValidator(now),
		jwt.ExpirationTimeValidator(now))
	header, err := jwt.Verify([]byte(token), resolver, &payload, validator)
	if err != nil {
		return nil, err
	}
	sub := permission.Subject_NONE
	if i, err := strconv.Atoi(payload.Subject); err != nil {
		sub = permission.Subject(permission.Subject_value[payload.Subject])
	} else {
		sub = permission.Subject(i)
	}
	return &secret.Info{
		Type:      secret.Type(secret.Type_value[header.ContentType]),
		Algorithm: alg,
		Subject:   sub,
		Issuer:    payload.Issuer,
		Account:   accountID,
		KeyId:     keyID,
		Roles:     payload.Audience,
		IssuedAt:  datetime.WithTime(payload.IssuedAt.Time).Proto(),
		ExpiredAt: datetime.WithTime(payload.ExpirationTime.Time).Proto(),
	}, nil
}
