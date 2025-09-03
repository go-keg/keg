package helpers

import (
	"crypto/md5" //nolint:gosec
	"crypto/sha256"
	"encoding/hex"
	"time"
)

func If[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}

// Throttle 节流函数
func Throttle(ttl time.Duration) func(fn func()) {
	var lastExecAt time.Time
	return func(fn func()) {
		if time.Now().After(lastExecAt.Add(ttl)) {
			fn()
			lastExecAt = time.Now()
		}
	}
}

func SHA256(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func MD5(data string) string {
	hash := md5.Sum([]byte(data)) //nolint:gosec
	return hex.EncodeToString(hash[:])
}
