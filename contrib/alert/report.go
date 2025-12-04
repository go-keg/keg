package alert

import (
	"context"
	"time"

	"github.com/go-keg/keg/contrib/helpers"
)

type Report func(ctx context.Context, record *ErrorRecord) (string, bool)

// DefaultReport 默认上报策略，只上报首次类似错误
var DefaultReport Report = func(ctx context.Context, record *ErrorRecord) (string, bool) {
	return record.RawMsg, record.Count == 1
}

// FibReport 按Fibonacci数列间隔小时上报，1,2,3,5,8,13,21,34,55 小时
var FibReport Report = func(ctx context.Context, record *ErrorRecord) (string, bool) {
	diff := int(time.Now().Sub(record.FirstSeen) / time.Hour)
	fib := helpers.Fib(record.AlertCount + 2)
	if diff > fib || record.Count == 1 {
		record.AlertCount++
		return record.RawMsg, true
	}
	return "", false
}
