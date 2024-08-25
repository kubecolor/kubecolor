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
	colorPrefix  = color.MustParseColor("gray:italic")
	colorDebug   = color.MustParseColor("gray:italic")
	colorInfo    = color.MustParseColor("green")
	colorWarn    = color.MustParseColor("yellow")
	colorError   = color.MustParseColor("red")
	colorMessage = color.MustParseColor("white")
	colorKey     = color.MustParseColor("cyan")
	colorValue   = color.MustParseColor("yellow")
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
	return h.Level >= level
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
		buf.WriteString(colorDebug.Render("[debug]"))
		buf.WriteByte(' ')
	case slog.LevelError:
		buf.WriteString(colorError.Render("[error]"))
		buf.WriteByte(' ')
	case slog.LevelInfo:
		buf.WriteByte(' ')
		buf.WriteString(colorInfo.Render("[info]"))
		buf.WriteByte(' ')
	case slog.LevelWarn:
		buf.WriteByte(' ')
		buf.WriteString(colorWarn.Render("[warn]"))
		buf.WriteByte(' ')
	default:
		fmt.Fprintf(&buf, "[Level(%d)] ", record.Level)
	}

	if record.Message != "" {
		buf.WriteString(colorMessage.Render(record.Message))
		buf.WriteByte(' ')
	}

	record.Attrs(func(attr slog.Attr) bool {
		switch attr.Value.Kind() {
		case slog.KindDuration:
			fmt.Fprintf(&buf, "%s=%s",
				colorKey.Render(attr.Key),
				colorValue.Render(duration.HumanDuration(attr.Value.Duration())))
		default:
			s := attr.Value.String()
			if needsQuoting(s) {
				fmt.Fprintf(&buf, "%s=%s", colorKey.Render(attr.Key), colorValue.Render(s))
			} else {
				fmt.Fprintf(&buf, "%s=%q", colorKey.Render(attr.Key), colorValue.Render(s))
			}
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
		if !unicode.IsPrint(r) {
			return true
		}
	}
	return false
}
