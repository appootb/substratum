package auth

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/appootb/protobuf/go/permission"
	"github.com/appootb/substratum/secret"
	"github.com/appootb/substratum/util/datetime"
	"github.com/appootb/substratum/util/hash"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/gbrlsnchs/jwt/v3/jwtutil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type JwtToken struct{}

func NewJwtToken() Token {
	return &JwtToken{}
}

func (t *JwtToken) sign(s *permission.Secret, seed []byte) (string, error) {
	issuedAt := datetime.FromProtoTime(s.GetIssuedAt()).Time
	payload := &jwt.Payload{
		Issuer:         s.GetIssuer(),
		Subject:        s.GetSubject(),
		Audience:       s.GetRoles(),
		ExpirationTime: jwt.NumericDate(datetime.FromProtoTime(s.GetExpiredAt()).Time),
		NotBefore:      jwt.NumericDate(issuedAt.Add(-time.Minute)),
		IssuedAt:       jwt.NumericDate(issuedAt),
		JWTID:          strconv.FormatUint(s.GetAccountId(), 10),
	}
	keyID := jwt.KeyID(fmt.Sprintf("%v-%v", s.GetAccountId(), hash.Sum(s.GetKeyId())))
	token, err := jwt.Sign(payload, jwt.NewHS256(seed), keyID)
	if err != nil {
		return "", err
	}
	return string(token), nil
}

func (t *JwtToken) TokenLevelValidator(level permission.TokenLevel) jwt.Validator {
	return func(pl *jwt.Payload) error {
		if level == permission.TokenLevel_NONE_TOKEN {
			return nil
		}
		for _, v := range pl.Audience {
			if !strings.HasSuffix(v, "_TOKEN") {
				continue
			}
			if permission.TokenLevel_value[v] < int32(level) {
				return status.Error(codes.PermissionDenied, fmt.Sprintf("expect: %v, actual: %v", level, v))
			}
			return nil
		}
		return jwt.ErrAudValidation
	}
}

// Generate a new token with specified options.
func (t *JwtToken) Generate(s *permission.Secret) (string, error) {
	seed, err := secret.DefaultSeed.New(s.GetAccountId(), hash.Sum(s.GetKeyId()))
	if err != nil {
		return "", err
	}
	return t.sign(s, seed)
}

// Refresh the token with expired time renewed.
func (t *JwtToken) Refresh(s *permission.Secret) (string, error) {
	seed, err := secret.DefaultSeed.Get(s.GetAccountId(), hash.Sum(s.GetKeyId()))
	if err != nil {
		return "", err
	}
	return t.sign(s, seed)
}

// Parse and verify the token string.
func (t *JwtToken) Verify(token string, level permission.TokenLevel) (*permission.Secret, error) {
	var (
		accountID uint64
		keyID     string
	)
	resolver := &jwtutil.Resolver{
		New: func(header jwt.Header) (jwt.Algorithm, error) {
			keyIDs := strings.Split(header.KeyID, "-")
			if len(keyIDs) != 2 {
				return nil, jwt.ErrAlgValidation
			}
			keyID = keyIDs[1]
			accountID, _ = strconv.ParseUint(keyIDs[0], 10, 64)
			seed, err := secret.DefaultSeed.Get(accountID, keyID)
			if err != nil {
				return nil, err
			}
			return jwt.NewHS256(seed), nil
		},
	}
	// Verify
	now := time.Now()
	payload := jwt.Payload{}
	validator := jwt.ValidatePayload(&payload,
		t.TokenLevelValidator(level),
		jwt.NotBeforeValidator(now),
		jwt.ExpirationTimeValidator(now))
	_, err := jwt.Verify([]byte(token), resolver, &payload, validator)
	if err != nil {
		return nil, err
	}

	return &permission.Secret{
		Issuer:    &payload.Issuer,
		Subject:   &payload.Subject,
		AccountId: &accountID,
		KeyId:     &keyID,
		Roles:     payload.Audience,
		Metadata:  map[string]string{},
		IssuedAt:  datetime.WithTime(payload.IssuedAt.Time).Proto(),
		ExpiredAt: datetime.WithTime(payload.ExpirationTime.Time).Proto(),
	}, nil
}

// Revoke the token signed to the specified account and device.
func (t *JwtToken) Revoke(accountID uint64, keyID string) error {
	return secret.DefaultSeed.Revoke(accountID, hash.Sum(keyID))
}

// Revoke all tokens signed to the specified account.
func (t *JwtToken) RevokeAll(accountID uint64) error {
	return secret.DefaultSeed.RevokeAll(accountID)
}
