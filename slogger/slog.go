package slogger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/jsasg/gopkg/slogger/handlers"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

type cfg struct {
	writer io.Writer
	rtime  time.Duration // 日志分隔时间
	mage   time.Duration // 日志保存最长时间
	source bool          // 是否记录日志位置
	level  slog.Level    // 日志级别
	format string        // 日志格式
}

type Option func(c *cfg)

func (opt Option) apply(c *cfg) {
	opt(c)
}

// WithPath 设置日志存储路径
func WithPath(path string) Option {
	return func(c *cfg) {
		var logpath = "runtime/logs/"
		if path != "" {
			logpath = path
		}
		if fidx := strings.LastIndex(path, "/"); fidx == -1 || fidx != len(path)-1 {
			logpath = fmt.Sprintf("%s%s", path, "/")
		}
		writer, _ := rotatelogs.New(
			fmt.Sprintf("%s%s", logpath, "/%Y%m/%d.log"),
			rotatelogs.WithClock(rotatelogs.Local),
			rotatelogs.WithRotationTime(c.rtime),
			rotatelogs.WithMaxAge(c.mage),
		)
		c.writer = io.MultiWriter(os.Stdout, writer)
	}
}

// WithRotationTime 设置日志分隔时间
func WithRotationTime(rotationTime time.Duration) Option {
	return func(c *cfg) {
		c.rtime = rotationTime
	}
}

// WithRetentionDays 设置日志保留最长天数
func WithRetentionDays(days int64) Option {
	return func(c *cfg) {
		c.mage = time.Hour * 2 * time.Duration(days)
	}
}

// WithMaxAge 设置日志保存最长时间
func WithMaxAge(maxAge time.Duration) Option {
	return func(c *cfg) {
		c.mage = maxAge
	}
}

// WithSource 设置是否记录日志位置
func WithSource(source bool) Option {
	return func(c *cfg) {
		c.source = source
	}
}

// WithLevel 设置日志记录级别
func WithLevel(level string) Option {
	return func(c *cfg) {
		switch level {
		case "debug":
			c.level = slog.LevelDebug
		case "info":
			c.level = slog.LevelInfo
		case "warn":
			c.level = slog.LevelWarn
		case "error":
			c.level = slog.LevelError
		default:
			c.level = slog.LevelError
		}
	}
}

// WithFormat 设置日志记录格式
func WithFormat(format string) Option {
	return func(c *cfg) {
		c.format = format
	}
}

// NewSlog 设置日志路径
func NewSlog(opts ...Option) *slog.Logger {
	var (
		cfg = &cfg{
			writer: os.Stdout,
			rtime:  time.Hour * 24,
			mage:   time.Hour * 24 * 7,
			source: true,
			level:  slog.LevelError,
		}
		handler slog.Handler
	)
	if len(opts) > 0 {
		for _, opt := range opts {
			opt.apply(cfg)
		}
	}
	handlerOpts := &slog.HandlerOptions{
		AddSource: cfg.source,
		Level:     cfg.level,
	}
	handler = handlers.NewTextHandler(cfg.writer, handlerOpts)
	if cfg.format == "json" {
		handler = slog.NewJSONHandler(cfg.writer, handlerOpts)
	}
	logger := slog.New(handler)
	return logger
}

// Default 设置默认logger
func Default(logger *slog.Logger) {
	slog.SetDefault(logger)
}
