package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const source = "abcdefghijklmnopqrstuvwxyz"

func init() {
	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())
}

// 生成 min到max之间的随机整数
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// 生成指定长度的随机字符串
func RandomString(n int) string {
	var sb strings.Builder
	k := len(source)

	for i := 0; i < n; i++ {
		c := source[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{USD, CNY}
	n := len(currencies)

	return currencies[rand.Intn(n)]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
