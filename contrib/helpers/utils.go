package helpers

import (
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
