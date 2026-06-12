package log

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type DailyTaggedWriter struct {
	baseDir string
	writers map[string]*taggedWriter
	mu      sync.Mutex
}

type taggedWriter struct {
	currDate string
	file     *os.File
}

func NewDailyTaggedWriter(baseDir string) *DailyTaggedWriter {
	return &DailyTaggedWriter{
		baseDir: baseDir,
		writers: make(map[string]*taggedWriter),
	}
}

func (d *DailyTaggedWriter) WriteWithTag(tag string, p []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	date := time.Now().Format("2006-01-02")

	writer, ok := d.writers[tag]
	if !ok || writer.currDate != date {
		// 新建或日期变更，重新打开文件
		if ok && writer.file != nil {
			_ = writer.file.Close()
		}
		filename := filepath.Join(d.baseDir, fmt.Sprintf("%s-%s.log", tag, date))
		_ = os.MkdirAll(filepath.Dir(filename), 0755)
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		writer = &taggedWriter{
			currDate: date,
			file:     file,
		}
		d.writers[tag] = writer
	}

	_, err := writer.file.Write(p)
	return err
}

func (d *DailyTaggedWriter) Sync() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, w := range d.writers {
		if w.file != nil {
			_ = w.file.Sync()
		}
	}
	return nil
}

type TaggedRoutingCore struct {
	defaultName  string
	encoder      zapcore.Encoder
	writer       *DailyTaggedWriter
	levelEnabler zapcore.LevelEnabler
}

func NewTaggedRoutingCore(defaultName string, enc zapcore.Encoder, writer *DailyTaggedWriter, level zapcore.LevelEnabler) zapcore.Core {
	return &TaggedRoutingCore{
		defaultName:  defaultName,
		encoder:      enc,
		writer:       writer,
		levelEnabler: level,
	}
}

func (c *TaggedRoutingCore) Enabled(lvl zapcore.Level) bool {
	return c.levelEnabler.Enabled(lvl)
}

func (c *TaggedRoutingCore) With(fields []zapcore.Field) zapcore.Core {
	clone := *c
	return &clone
}

func (c *TaggedRoutingCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *TaggedRoutingCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	tag := c.defaultName
	for _, f := range fields {
		if f.Key == "tag" && f.Type == zapcore.StringType {
			tag = f.String
			break
		}
	}

	buf, err := c.encoder.EncodeEntry(ent, fields)
	if err != nil {
		return err
	}
	return c.writer.WriteWithTag(tag, buf.Bytes())
}

func (c *TaggedRoutingCore) Sync() error {
	return c.writer.Sync()
}

/**
kratos
*/

type ZapTaggedLogger struct {
	zap *zap.Logger
}

func NewZapTaggedLogger(z *zap.Logger) *ZapTaggedLogger {
	return &ZapTaggedLogger{zap: z}
}

type Channel string

const ChannelName Channel = "channel"

func (l *ZapTaggedLogger) Log(level log.Level, keyvals ...any) error {
	if len(keyvals)%2 != 0 {
		return fmt.Errorf("invalid keyvals, must be even")
	}

	var fields []zap.Field
	var msg string
	for i := 0; i < len(keyvals); i += 2 {
		k := fmt.Sprint(keyvals[i])
		v := keyvals[i+1]
		if k == "msg" {
			msg = fmt.Sprint(v)
		} else {
			fields = append(fields, zap.Any(k, v))
		}
	}

	logger := l.zap
	switch level {
	case log.LevelDebug:
		logger.Debug(msg, fields...)
	case log.LevelInfo:
		logger.Info(msg, fields...)
	case log.LevelWarn:
		logger.Warn(msg, fields...)
	case log.LevelError:
		logger.Error(msg, fields...)
	default:
		logger.Info(msg, fields...)
	}
	return nil
}

func NewZapWithTaggedRouting(options *ZapOptions) *zap.Logger {
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

	// core := NewTaggedRoutingCore(defaultName, encoder, writer, zapcore.DebugLevel)
	// 判断日志等级的interface
	level := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= options.Level
	})
	core := zapcore.NewTee(NewTaggedRoutingCore(
		options.Filename,
		encoder,
		NewDailyTaggedWriter(options.Filename),
		level,
	))

	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))
}
