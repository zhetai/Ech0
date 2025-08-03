package util

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5Encrypt 对内容进行 MD5 编码
func MD5Encrypt(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))
	hashInBytes := hash.Sum(nil)
	return hex.EncodeToString(hashInBytes)
}
