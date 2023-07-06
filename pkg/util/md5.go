package util

import (
	"crypto/md5"
	"encoding/hex"
)

// 计算 MD5
func EncodeMd5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}