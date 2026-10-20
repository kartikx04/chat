// internal/logger/slog.go
package logger

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type storedAttr struct {
	attr   slog.Attr
	groups []string
}

type colorHandler struct {
	mu     *sync.Mutex
	out    io.Writer
	lvl    slog.Level
	attrs  []storedAttr
	groups []string
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

func newColorHandler(out io.Writer, level slog.Level) *colorHandler {
	return &colorHandler{
		mu:  &sync.Mutex{},
		out: out,
		lvl: level,
	}
}

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

func (h *colorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}

	clone := *h
	clone.attrs = append([]storedAttr(nil), h.attrs...)

	for _, attr := range attrs {
		clone.attrs = append(clone.attrs, storedAttr{
			attr:   attr,
			groups: append([]string(nil), h.groups...),
		})
	}
	return &clone

}
func (h *colorHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	clone := *h
	clone.groups = append([]string(nil), h.groups...)
	clone.groups = append(clone.groups, name)

	return &clone
}

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

	// message — bold white
	buf.WriteString(bold)
	buf.WriteString(r.Message)
	buf.WriteString(reset)
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

	for _, stored := range h.attrs {
		writeAttr(&buf, stored.attr, stored.groups)
	}

	// key=value attrs — key in blue, value plain
	r.Attrs(func(attr slog.Attr) bool {
		writeAttr(&buf, attr, h.groups)
		return true
	})

	buf.WriteByte('\n')

	h.mu.Lock()
	_, err := h.out.Write(buf.Bytes())
	h.mu.Unlock()
	return err
}

func writeAttr(buf *bytes.Buffer, attr slog.Attr, groups []string) {
	attr.Value = attr.Value.Resolve()

	if attr.Value.Kind() == slog.KindGroup {
		nextGroups := append([]string(nil), groups...)
		if attr.Key != "" {
			nextGroups = append(nextGroups, attr.Key)
		}

		for _, child := range attr.Value.Group() {
			writeAttr(buf, child, nextGroups)
		}

		return
	}

	if attr.Key == "" {
		return
	}

	key := strings.Join(groups, ".")
	if key != "" {
		key += "."
	}
	key += attr.Key

	value := fmt.Sprint(attr.Value.Any())

	buf.WriteByte(' ')
	buf.WriteString(blue)
	buf.WriteString(key)
	buf.WriteString(reset)
	buf.WriteByte('=')
	buf.WriteString(cyan)

	if containsSpace(value) || attr.Value.Kind() == slog.KindString {
		buf.WriteByte('"')
		buf.WriteString(value)
		buf.WriteByte('"')
	} else {
		buf.WriteString(value)
	}

	buf.WriteString(reset)
}

func containsSpace(s string) bool {
	for _, c := range s {
		if c == ' ' || c == ':' || c == '=' {
			return true
		}
	}
	return false
}

func Init(env string) *slog.Logger {
	var handler slog.Handler

	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: false,
		})
	} else {
		handler = newColorHandler(os.Stdout, slog.LevelDebug)
	}

	log := slog.New(handler)
	slog.SetDefault(log)
	return log
}
