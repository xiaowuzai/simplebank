package token

import "time"

// Maker 是用来管理 token 的接口
type Maker interface {
	// CreateToken 根据用户名和duration 创建 token
	CreateToken(usernames string, role string, duration time.Duration) (string, *Payload, error)

	// VerifyToken 验证 token
	VerifyToken(token string) (*Payload, error)
}
