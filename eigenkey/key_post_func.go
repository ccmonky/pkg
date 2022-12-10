package eigenkey

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

var (
	keyPostFuncRegistry = map[string]KeyPostFunc{
		"md5":      MD5,    // 长度32
		"sha1":     SHA1,   // 长度40
		"sha256":   SHA256, // 长度64
		"prefix64": Prefix64,
	}
)

// KeyPostFunc 对key进行后处理
type KeyPostFunc func(string) string

// MD5 计算md5
func MD5(key string) string {
	h := md5.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}

// SHA1 计算SHA1
func SHA1(key string) string {
	h := sha1.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}

// SHA256 计算SHA256
func SHA256(key string) string {
	h := sha256.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}

// Prefix64 取前64个字符
func Prefix64(key string) string {
	if len(key) > 64 {
		return key[:64]
	}
	return key
}
