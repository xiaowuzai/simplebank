package token

import (
	"encoding/json"
	"time"

	"aidanwoods.dev/go-paseto"
)

// PasetoMaker 生成 paseto token
type PasetoMaker struct {
	symmetric paseto.V4SymmetricKey
}

// NewPasetoMaker  Symmetric
func NewPasetoMaker(hexKey string) (Maker, error) {
	symmetric, err := paseto.V4SymmetricKeyFromHex(hexKey)
	if err != nil {
		return nil, err
	}

	maker := &PasetoMaker{
		symmetric: symmetric,
	}

	return maker, nil
}

// CreateToken 根据用户名和duration 创建 token
func (m *PasetoMaker) CreateToken(usernames string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(usernames, duration)
	if err != nil {
		return "", payload, nil
	}

	jsonClaim, err := json.Marshal(payload)
	if err != nil {
		return "", nil, err
	}

	token, err := paseto.NewTokenFromClaimsJSON(jsonClaim, []byte{})
	if err != nil {
		return "", nil, err
	}
	token.SetIssuedAt(payload.IssuedAt)
	token.SetExpiration(payload.ExpiredAt)

	// key := paseto.NewV4SymmetricKey() // don't share this!!
	encrypted := token.V4Encrypt(m.symmetric, nil)
	return encrypted, payload, nil
}

// VerifyToken 验证 token
func (m *PasetoMaker) VerifyToken(signed string) (*Payload, error) {
	parser := paseto.NewParser()

	rules := []paseto.Rule{
		paseto.NotExpired(),
	}
	parser.SetRules(rules)
	token, err := parser.ParseV4Local(m.symmetric, signed, nil)
	if err != nil {
		return nil, ErrExpiredToken
	}

	payload := &Payload{}
	if err := json.Unmarshal(token.ClaimsJSON(), payload); err != nil {
		return nil, ErrInvalidToken
	}

	return payload, nil
}

// NewPasetoSymmetricKey 生成对称密钥字符串
// func NewPasetoSymmetricKey() string {
// 	return paseto.NewV4SymmetricKey().ExportHex()
// }
