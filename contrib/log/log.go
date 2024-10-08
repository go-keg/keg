package log

import (
	"fmt"
	"github.com/go-keg/keg/contrib/config"
	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"path"
	"time"
)

type KeyValue func() (string, any)

func ServiceID(val string) KeyValue {
	return func() (string, any) {
		return "service.id", val
	}
}

func ServiceName(val string) KeyValue {
	return func() (string, any) {
		return "service.name", val
	}
}

func ServiceVersion(val string) KeyValue {
	return func() (string, any) {
		return "service.version", val
	}
}

func NewLoggerFromConfig(conf config.Log, name string, keyValues ...KeyValue) log.Logger {
	options := &ZapOptions{
		Filename:     path.Join(conf.Dir, name),
		Level:        Level(conf.Level),
		MaxAge:       time.Duration(conf.MaxAge) * time.Hour * 24,
		RotationTime: time.Duration(conf.RotationTime) * time.Hour * 24,
	}
	var values []any
	values = append(values, "ts", log.DefaultTimestamp)
	values = append(values, "caller", log.DefaultCaller)
	for _, keyValue := range keyValues {
		k, v := keyValue()
		values = append(values, k, v)
	}
	return log.With(NewLogger(options), values...)
}

func NewLogger(options *ZapOptions) log.Logger {
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		LevelKey:    "level",
		LineEnding:  zapcore.DefaultLineEnding,
		EncodeLevel: zapcore.LowercaseLevelEncoder,
	})

	level := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= options.Level
	})

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(writer(options)), level),
	)

	return &zapLogger{
		logger: zap.New(core),
	}
}

type zapLogger struct {
	logger *zap.Logger
}

func (l *zapLogger) Log(level log.Level, keyValues ...interface{}) error {
	if len(keyValues) == 0 || len(keyValues)%2 != 0 {
		l.logger.Warn(fmt.Sprint("keyValues must appear in pairs: ", keyValues))
		return nil
	}

	var data []zap.Field
	for i := 0; i < len(keyValues); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyValues[i]), fmt.Sprint(keyValues[i+1])))
	}
	switch level {
	case log.LevelDebug:
		l.logger.Debug("", data...)
	case log.LevelInfo:
		l.logger.Info("", data...)
	case log.LevelWarn:
		l.logger.Warn("", data...)
	case log.LevelError:
		l.logger.Error("", data...)
	case log.LevelFatal:
		l.logger.Fatal("", data...)
	}
	return nil
}
