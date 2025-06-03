package log

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// Option log config options
type Option = logger.Config

// Level log level
type Level = logger.LogLevel

// Log Level Value
const (
	Silent Level = iota + 1
	Error
	Warn
	Info
)

// Ormlogger records ORM logs with a structured logger
type Ormlogger struct {
	*slog.Logger
	option logger.Config
}

// DefaultOrmlogger create a default Ormlogger instance
func DefaultOrmlogger() *Ormlogger {
	return NewOrmlogger(nil, Option{
		ParameterizedQueries: true,
		LogLevel:             Info,
		SlowThreshold:        time.Millisecond * 300,
	})
}

// NewOrmlogger creates a new Ormlogger instance
func NewOrmlogger(logger *slog.Logger, opts Option) *Ormlogger {
	if logger == nil {
		logger = slog.Default()
	}
	return &Ormlogger{Logger: logger, option: opts}
}

// LogMode log mode
func (o *Ormlogger) LogMode(level logger.LogLevel) logger.Interface {
	o.option.LogLevel = level
	return NewOrmlogger(o.Logger, o.option)
}

// Info prints info message
func (o *Ormlogger) Info(ctx context.Context, msg string, data ...any) {
	if o.option.LogLevel >= logger.Info {
		o.Logger.InfoContext(ctx, msg, slog.Any("data", data))
	}
}

// Warn prints warn message
func (o *Ormlogger) Warn(ctx context.Context, msg string, data ...any) {
	if o.option.LogLevel >= logger.Warn {
		o.Logger.WarnContext(ctx, msg, slog.Any("data", data))
	}
}

// Error prints error message
func (o *Ormlogger) Error(ctx context.Context, msg string, data ...any) {
	if o.option.LogLevel >= logger.Error {
		o.Logger.ErrorContext(ctx, msg, slog.Any("data", data))
	}
}

// Trace logs SQL queries and execution details
func (o *Ormlogger) Trace(ctx context.Context, begin time.Time,
	fc func() (sql string, rowsAffected int64), err error) {
	if o.option.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	attrs := []any{
		"caller", shortCaller(utils.FileWithLineNum()),
		"sql", sql,
		"rows", rows,
		"elapsed", elapsed,
	}

	switch {
	case err != nil && o.option.LogLevel >= logger.Error:
		if errors.Is(err, gorm.ErrRecordNotFound) &&
			o.option.IgnoreRecordNotFoundError {
			return
		}
		o.Logger.ErrorContext(ctx, "Query Error", append(attrs, "error", err)...)

	case elapsed > o.option.SlowThreshold && o.option.SlowThreshold > 0 &&
		o.option.LogLevel >= logger.Warn:
		msg := fmt.Sprintf("SLOW SQL >= %v", o.option.SlowThreshold)
		o.Logger.WarnContext(ctx, msg, attrs...)

	case o.option.LogLevel >= logger.Info:
		o.Logger.InfoContext(ctx, "Query Executed", attrs...)
	}
}

func shortCaller(fullPath string) string {
	parts := strings.Split(fullPath, "/")
	if len(parts) < 2 {
		return fullPath
	}
	dir := parts[len(parts)-2]
	fileWithLine := parts[len(parts)-1]
	return filepath.Join(dir, fileWithLine)
}
