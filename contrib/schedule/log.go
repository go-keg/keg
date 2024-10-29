package schedule

import "github.com/go-kratos/kratos/v2/log"

type cronLogger struct {
	log *log.Helper
}

func (l cronLogger) Debug(msg string, keysAndValues ...any) {
	keysAndValues = append([]any{"msg", msg}, keysAndValues...)
	l.log.Debugw(keysAndValues...)
}

func (l cronLogger) Info(msg string, keysAndValues ...any) {
	keysAndValues = append([]any{"msg", msg}, keysAndValues...)
	l.log.Infow(keysAndValues...)
}

func (l cronLogger) Error(err error, msg string, keysAndValues ...any) {
	keysAndValues = append([]any{"msg", msg, "err", err}, keysAndValues...)
	l.log.Errorw(keysAndValues...)
}
