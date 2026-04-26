package logger

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"gorm.io/gorm/logger"
)

type SlogGormLogger struct {
	SlowThreshold time.Duration
	env           string
}

func NewSlogGormLogger(env string) *SlogGormLogger {
	return &SlogGormLogger{
		SlowThreshold: 200 * time.Millisecond,
		env:           env,
	}
}

func (l *SlogGormLogger) LogMode(_ logger.LogLevel) logger.Interface {
	return l
}

func (l *SlogGormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	slog.InfoContext(ctx, fmt.Sprintf(msg, args...))
}

func (l *SlogGormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	slog.WarnContext(ctx, fmt.Sprintf(msg, args...))
}

func (l *SlogGormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	slog.ErrorContext(ctx, fmt.Sprintf(msg, args...))
}

func (l *SlogGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil {
		slog.ErrorContext(ctx, "db error",
			"error", err,
			"rows", rows,
			"duration_ms", elapsed.Milliseconds(),
		)
		return
	}

	// Skip schema introspection queries entirely (AutoMigrate noise)
	if strings.Contains(sql, "information_schema") ||
		strings.Contains(sql, "pg_catalog") ||
		strings.Contains(sql, "pg_attribute") {
		return
	}

	if l.env == "development" {
		slog.DebugContext(ctx, "db query",
			"sql", sql,
			"rows", rows,
			"duration_ms", elapsed.Milliseconds(),
		)
		return
	}

	if elapsed > l.SlowThreshold {
		slog.WarnContext(ctx, "slow db query",
			"sql", sql,
			"rows", rows,
			"duration_ms", elapsed.Milliseconds(),
			"threshold_ms", l.SlowThreshold.Milliseconds(),
		)
	}
}
