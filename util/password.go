package util

import "golang.org/x/crypto/bcrypt"

// GeneratePassword 通过 bcrypt 加密字符串
func HashPassword(password string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(pass), nil
}

// CheckPassword 验证输入的密码和加密之后的密码是否匹配
func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
