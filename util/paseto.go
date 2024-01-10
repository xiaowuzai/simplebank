package util

import "aidanwoods.dev/go-paseto"

// NewPasetoSymmetricKey 生成对称密钥字符串
func NewPasetoSymmetricKey() string {
	return paseto.NewV4SymmetricKey().ExportHex()
}
