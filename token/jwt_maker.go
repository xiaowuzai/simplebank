package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

// JWTMaker 是一个 JSON Web Token Maker
type JWTMaker struct {
	secretKey string
}

type jwtClaimPayload struct {
	jwt.RegisteredClaims
}

// NewJWTMaker 创建一个 JWTMaker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

// CreateToken 根据用户名和duration 创建 token
func (m *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil, err
	}

	claims := &jwtClaimPayload{
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(payload.ExpiredAt),
			IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
			Issuer:    payload.Username,
			ID:        payload.ID,
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", nil, err
	}
	return token, payload, nil
}

// VerifyToken 验证 token
func (m *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(m.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &jwtClaimPayload{}, keyFunc)
	if err != nil {
		return nil, err
	}
	claims, ok := jwtToken.Claims.(*jwtClaimPayload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return &Payload{
		ID:        claims.ID,
		Username:  claims.Issuer,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiredAt: claims.ExpiresAt.Time,
	}, nil

}
