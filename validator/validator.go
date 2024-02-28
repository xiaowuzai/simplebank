package validator

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

// 验证字符串长度
func ValidateString(value string, min, max int) error {
	n := len(value)
	if n < min || n > max {
		return fmt.Errorf("长度限制不符合规范%d-%d", min, max)
	}
	return nil
}

// 验证用户名
func ValidateUsername(value string) error {
	if err := ValidateString(value, 3, 30); err != nil {
		return err
	}

	if !isValidUsername(value) {
		return fmt.Errorf("只能包含: 大小写字母、数字、下划线")
	}
	return nil
}

// 验证密码
func ValidatePassword(value string) error {
	if err := ValidateString(value, 6, 30); err != nil {
		return err
	}
	return nil
}

// 验证邮箱
func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 200); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("邮箱地址错误")
	}
	return nil
}

// 验证 全名
func ValidateFullName(value string) error {
	if err := ValidateString(value, 3, 30); err != nil {
		return err
	}

	if !isValidFullName(value) {
		return fmt.Errorf("只能包含: 大小写字母、空格")
	}
	return nil
}

// 验证 emailId
func ValidateEmailId(value int64) error {
	if value <= 0 {
		return fmt.Errorf("只能是正整数")
	}
	return nil
}

// 验证 Email 的验证码长度
func ValidateEmailSecretCode(value string) error {
	return ValidateString(value, 32, 32)
}
