package testutil

import (
	"io"
	"log/slog"
	"testing"
)

func SetTestLogger(t testing.TB, w io.Writer) {
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				return slog.Attr{}
			}
			return a
		},
	})))
	t.Cleanup(func() {
		slog.SetDefault(previous)
	})
}
