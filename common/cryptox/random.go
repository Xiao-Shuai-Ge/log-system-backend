package cryptox

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomString 生成一个加密安全的指定长度（字节）的随机字符串（返回 2*n 长度的十六进制字符串）
func GenerateRandomString(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
