package helpers

import (
	"crypto/md5" //nolint:gosec
	"crypto/sha256"
	"fmt"
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

func HashWithSHA256(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func HashWithMD5(data string) string {
	hash := md5.Sum([]byte(data)) //nolint:gosec
	return fmt.Sprintf("%x", hash)
}
