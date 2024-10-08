package log

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"time"
)

type ZapOptions struct {
	Filename     string
	MaxAge       time.Duration
	RotationTime time.Duration
	Fields       map[string]string
	Level        zapcore.Level
}

func NewZapLog(options *ZapOptions) *zap.Logger {
	// 设置一些基本日志格式
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		TimeKey:       "ts",
		StacktraceKey: "trace",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(time.RFC3339))
		},
		EncodeCaller: zapcore.ShortCallerEncoder,
		LineEnding:   zapcore.DefaultLineEnding,
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
	})

	// 添加自定义字段
	for field, value := range options.Fields {
		encoder.AddString(field, value)
	}

	// 判断日志等级的interface
	level := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= options.Level
	})

	// 获取 info、error日志文件的io.Writer 抽象 getWriter() 在下方实现
	core := zapcore.NewTee(zapcore.NewCore(encoder, zapcore.AddSync(writer(options)), level))
	return zap.New(core, zap.AddCaller())
}

func writer(options *ZapOptions) io.Writer {
	hook, err := rotatelogs.New(
		options.Filename+"-%Y-%m-%d.log",
		rotatelogs.WithMaxAge(options.MaxAge),
		rotatelogs.WithRotationTime(options.RotationTime),
	)

	if err != nil {
		panic(err)
	}
	return hook
}

func Level(level string) zapcore.Level {
	levels := []string{"debug", "info", "warn", "error", "dpanic", "panic", "fatal"}
	for i, l := range levels {
		if l == level {
			return zapcore.Level(i - 1)
		}
	}
	panic(fmt.Sprintf("log: Invalid level string '%s'", level))
}
