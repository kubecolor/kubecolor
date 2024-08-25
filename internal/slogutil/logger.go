package slogutil

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"unicode"

	"github.com/kubecolor/kubecolor/config/color"
	"github.com/mattn/go-colorable"
	"k8s.io/apimachinery/pkg/util/duration"
)

var (
	colorPrefix  = color.MustParse("gray:italic")
	colorDebug   = color.MustParse("gray:italic")
	colorInfo    = color.MustParse("green")
	colorWarn    = color.MustParse("yellow")
	colorError   = color.MustParse("red")
	colorMessage = color.MustParse("white")
	colorKey     = color.MustParse("cyan")
	colorValue   = color.MustParse("light-yellow")
)

type SlogHandler struct {
	Writer io.Writer
	Level  slog.Level

	attrs []slog.Attr
}

var _ slog.Handler = &SlogHandler{}

type SlogHandlerOptions struct {
	Writer io.Writer
	Level  slog.Level
}

func NewSlogHandler(opt *SlogHandlerOptions) *SlogHandler {
	if opt == nil {
		opt = &SlogHandlerOptions{}
	}
	if opt.Writer == nil {
		opt.Writer = colorable.NewColorableStderr()
	}
	if opt.Level == 0 {
		opt.Level = slog.LevelError
	}
	return &SlogHandler{
		Writer: opt.Writer,
		Level:  opt.Level,
	}
}

// Enabled implements [slog.Handler].
func (h *SlogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return h.Level <= level
}

// WithAttrs implements [slog.Handler].
func (h *SlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	copy := *h
	copy.attrs = append(copy.attrs, attrs...)
	return &copy
}

// WithGroup implements [slog.Handler].
func (h *SlogHandler) WithGroup(name string) slog.Handler {
	panic("unimplemented")
}

// Handle implements [slog.Handler].
func (h *SlogHandler) Handle(_ context.Context, record slog.Record) error {
	var buf bytes.Buffer

	buf.WriteString(colorPrefix.Render("[kubecolor]"))
	buf.WriteByte(' ')
	switch record.Level {
	case slog.LevelDebug:
		buf.WriteString(colorDebug.Render("[DEBUG]"))
	case slog.LevelError:
		buf.WriteString(colorError.Render("[ERROR]"))
	case slog.LevelInfo:
		buf.WriteByte(' ')
		buf.WriteString(colorInfo.Render("[INFO]"))
	case slog.LevelWarn:
		buf.WriteByte(' ')
		buf.WriteString(colorWarn.Render("[WARN]"))
	default:
		fmt.Fprintf(&buf, "[LEVEL(%d)]", record.Level)
	}

	if record.Message != "" {
		buf.WriteByte(' ')
		buf.WriteString(colorMessage.Render(record.Message))
	}

	record.Attrs(func(attr slog.Attr) bool {
		buf.WriteByte(' ')
		switch attr.Value.Kind() {
		case slog.KindDuration:
			fmt.Fprintf(&buf, "%s=%s",
				colorKey.Render(attr.Key),
				colorValue.Render(duration.HumanDuration(attr.Value.Duration())))
		default:
			s := attr.Value.String()
			if needsQuoting(s) {
				s = fmt.Sprintf("%q", s)
			}
			fmt.Fprintf(&buf, "%s=%s", colorKey.Render(attr.Key), colorValue.Render(s))
		}
		return true
	})

	buf.WriteByte('\n')

	writer := h.Writer
	if writer == nil {
		writer = os.Stderr
	}
	_, err := h.Writer.Write(buf.Bytes())
	return err
}

func needsQuoting(s string) bool {
	if s == "" {
		return true
	}
	for _, r := range s {
		switch r {
		case ' ', '"', '\'':
			return true
		}
		if !unicode.IsPrint(r) {
			return true
		}
	}
	return false
}
