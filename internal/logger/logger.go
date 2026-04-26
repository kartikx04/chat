// internal/logger/logger.go
package logger

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type colorHandler struct {
	mu  sync.Mutex
	lvl slog.Level
}

// ANSI codes
const (
	reset  = "\033[0m"
	cyan   = "\033[36m"
	green  = "\033[32m"
	yellow = "\033[33m"
	red    = "\033[31m"
	gray   = "\033[90m"
	bold   = "\033[1m"
	blue   = "\033[34m"
)

func levelColor(l slog.Level) string {
	switch {
	case l < slog.LevelInfo:
		return cyan + "DBG" + reset
	case l < slog.LevelWarn:
		return green + "INF" + reset
	case l < slog.LevelError:
		return yellow + "WRN" + reset
	default:
		return red + "ERR" + reset
	}
}

func (h *colorHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.lvl
}

func (h *colorHandler) WithAttrs(attrs []slog.Attr) slog.Handler { return h }
func (h *colorHandler) WithGroup(name string) slog.Handler       { return h }

func (h *colorHandler) Handle(_ context.Context, r slog.Record) error {
	var buf bytes.Buffer

	// time — gray, short
	buf.WriteString(gray)
	buf.WriteString(r.Time.Format("15:04:05"))
	buf.WriteString(reset)
	buf.WriteByte(' ')

	// level — colored, 3 chars so columns align
	buf.WriteString(levelColor(r.Level))
	buf.WriteByte(' ')

	// source — gray, just filename:line
	if r.PC != 0 {
		frames := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := frames.Next()
		buf.WriteString(gray)
		buf.WriteString(filepath.Base(f.File))
		buf.WriteByte(':')
		buf.WriteString(fmt.Sprint(f.Line))
		buf.WriteString(reset)
		buf.WriteByte(' ')
	}

	// message — bold white
	buf.WriteString(bold)
	buf.WriteString(r.Message)
	buf.WriteString(reset)

	// key=value attrs — key in blue, value plain
	r.Attrs(func(a slog.Attr) bool {
		buf.WriteByte(' ')
		buf.WriteString(blue)
		buf.WriteString(a.Key)
		buf.WriteString(reset)
		buf.WriteByte('=')
		val := fmt.Sprintf("%v", a.Value.Any())
		// quote values with spaces so they're clearly bounded
		if containsSpace(val) {
			buf.WriteByte('"')
			buf.WriteString(val)
			buf.WriteByte('"')
		} else {
			buf.WriteString(val)
		}
		return true
	})

	buf.WriteByte('\n')

	h.mu.Lock()
	os.Stdout.Write(buf.Bytes())
	h.mu.Unlock()
	return nil
}

func containsSpace(s string) bool {
	for _, c := range s {
		if c == ' ' || c == ':' || c == '=' {
			return true
		}
	}
	return false
}

func Init(env string) {
	var handler slog.Handler

	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: false,
		})
	} else {
		handler = &colorHandler{lvl: slog.LevelDebug}
	}

	slog.SetDefault(slog.New(handler))
}
