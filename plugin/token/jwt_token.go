package token

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"strconv"
	"time"

	"github.com/appootb/jwt"
	"github.com/appootb/substratum/v2/credential"
	"github.com/appootb/substratum/v2/proto/go/common"
	"github.com/appootb/substratum/v2/proto/go/permission"
	"github.com/appootb/substratum/v2/proto/go/secret"
	"github.com/appootb/substratum/v2/token"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	// UnsupportedAlgorithm is returned when secret.Algorithm is not mapped to a JWT JWA.
	UnsupportedAlgorithm = errors.New("substratum: algorithm not supported")
	// ErrTokenTypeMismatch is returned when JWT cty does not match the request platform.
	ErrTokenTypeMismatch = errors.New("substratum: token type mismatch")
)

// Init registers JwtToken as the default token implementor when none is set.
func Init() {
	if token.Implementor() == nil {
		token.RegisterImplementor(&JwtToken{})
	}
}

// JwtToken implements token.Token using github.com/appootb/jwt.
type JwtToken struct{}

func (t *JwtToken) sign(s *secret.Info, key []byte) (string, error) {
	issuedAt := s.GetIssuedAt().AsTime().In(time.Local)
	claims := &jwt.RegisteredClaims{
		Issuer:         s.GetIssuer(),
		Subject:        s.GetSubject().String(),
		Audience:       jwt.Audience(s.GetRoles()),
		ExpirationTime: jwt.NewNumericDate(s.GetExpiredAt().AsTime().In(time.Local)),
		NotBefore:      jwt.NewNumericDate(issuedAt.Add(-time.Minute)),
		IssuedAt:       jwt.NewNumericDate(issuedAt),
	}
	alg, err := t.getAlgorithm(s.GetAlgorithm(), key)
	if err != nil {
		return "", err
	}
	val, err := jwt.Sign(claims, alg,
		jwt.WithContentType(s.GetType().String()),
		jwt.WithAccountID(s.GetAccount()),
		jwt.WithKeyID(strconv.FormatInt(s.GetKeyId(), 10)),
	)
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
		return jwt.NewRS512(priv, &priv.PublicKey), nil
	case secret.Algorithm_PSS:
		priv, err := x509.ParsePKCS1PrivateKey(key)
		if err != nil {
			return nil, err
		}
		return jwt.NewPS512(priv, &priv.PublicKey), nil
	case secret.Algorithm_ECDSA:
		priv, err := x509.ParseECPrivateKey(key)
		if err != nil {
			return nil, err
		}
		return jwt.NewES512(priv, &priv.PublicKey), nil
	case secret.Algorithm_EdDSA:
		priv := ed25519.PrivateKey(key)
		return jwt.NewEdDSA(priv, priv.Public().(ed25519.PublicKey)), nil
	default:
		return nil, UnsupportedAlgorithm
	}
}

// NewSecretKey generates a new private key blob for alg.
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
	dur := s.GetExpiredAt().AsTime().Sub(s.GetIssuedAt().AsTime())
	if err = credential.ClientImplementor().Add(s.GetAccount(), s.GetKeyId(), key, dur); err != nil {
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
		dur := s.GetExpiredAt().AsTime().Sub(s.GetIssuedAt().AsTime())
		key, err = credential.ClientImplementor().Refresh(s.GetAccount(), s.GetKeyId(), dur)
	}
	if err != nil {
		return "", err
	}
	// Calculate new timestamp.
	issuedAt := s.GetIssuedAt().AsTime().In(time.Local)
	expiredAt := s.GetExpiredAt().AsTime().In(time.Local)
	now := time.Now()
	s.IssuedAt = timestamppb.New(now)
	s.ExpiredAt = timestamppb.New(now.Add(expiredAt.Sub(issuedAt)))
	return t.sign(s, key)
}

// Parse the metadata.
// After verification, JWT cty must match the request platform:
// PLATFORM_SERVER expects Type_SERVER, otherwise Type_CLIENT.
func (t *JwtToken) Parse(md *common.Metadata) (*secret.Info, error) {
	info, err := t.ParseRaw(md.GetToken())
	if err != nil {
		return nil, err
	}
	expected := secret.Type_CLIENT
	if md.GetPlatform()&common.Platform_PLATFORM_SERVER == common.Platform_PLATFORM_SERVER {
		expected = secret.Type_SERVER
	}
	if info.GetType() != expected {
		return nil, ErrTokenTypeMismatch
	}
	return info, nil
}

// ParseRaw parses a token string.
// Verification resolves the signing key from the JOSE header (account + kid + cty) first.
// Successful results are cached in-process for a short TTL (see SetParseCacheEnabled).
func (t *JwtToken) ParseRaw(raw string) (*secret.Info, error) {
	return cachedParseRaw(raw, func() (*secret.Info, error) {
		return t.parseRawUncached(raw)
	})
}

func (t *JwtToken) parseRawUncached(raw string) (*secret.Info, error) {
	var (
		accountID uint64
		keyID     int64
		alg       secret.Algorithm
	)
	now := time.Now()
	var claims jwt.RegisteredClaims
	header, err := jwt.VerifyWithKeyFunc([]byte(raw), func(h *jwt.Header) (jwt.Algorithm, error) {
		if h.KeyID == "" {
			return nil, jwt.ErrMalformed
		}
		parsedKeyID, parseErr := strconv.ParseInt(h.KeyID, 10, 64)
		if parseErr != nil {
			return nil, jwt.ErrMalformed
		}
		accountID = h.AccountID
		keyID = parsedKeyID

		var (
			key []byte
			err error
		)
		if h.ContentType == secret.Type_SERVER.String() {
			key, err = credential.ServerImplementor().Get(keyID)
		} else {
			key, err = credential.ClientImplementor().Get(accountID, keyID)
		}
		if err != nil {
			return nil, err
		}
		switch h.Algorithm {
		case "HS512":
			alg = secret.Algorithm_HMAC
		case "RS512":
			alg = secret.Algorithm_RSA
		case "PS512":
			alg = secret.Algorithm_PSS
		case "ES512":
			alg = secret.Algorithm_ECDSA
		case "EdDSA":
			alg = secret.Algorithm_EdDSA
		default:
			return nil, jwt.ErrInvalidAlgorithm
		}
		return t.getAlgorithm(alg, key)
	}, &claims,
		jwt.NotBeforeValidator(now, 0),
		jwt.ExpirationValidator(now, 0),
	)
	if err != nil {
		return nil, err
	}

	sub := permission.Subject_NONE
	if i, err := strconv.Atoi(claims.Subject); err != nil {
		sub = permission.Subject(permission.Subject_value[claims.Subject])
	} else {
		sub = permission.Subject(i)
	}

	var issuedAt, expiredAt time.Time
	if !claims.IssuedAt.IsZero() {
		issuedAt = claims.IssuedAt.Time
	}
	if !claims.ExpirationTime.IsZero() {
		expiredAt = claims.ExpirationTime.Time
	}

	return &secret.Info{
		Type:      secret.Type(secret.Type_value[header.ContentType]),
		Algorithm: alg,
		Subject:   sub,
		Issuer:    claims.Issuer,
		Account:   accountID,
		KeyId:     keyID,
		Roles:     claims.Audience,
		IssuedAt:  timestamppb.New(issuedAt),
		ExpiredAt: timestamppb.New(expiredAt),
	}, nil
}
