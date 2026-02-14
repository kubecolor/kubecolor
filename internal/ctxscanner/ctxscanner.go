// Package ctxscanner contains implementation of a scanner that takes in a
// [context.Context], and allows it to be paused.
package ctxscanner

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/kubecolor/kubecolor/internal/bytesutil"
)

// Scanner is a wrapper around [bufio.Scanner] that takes in a context.
type Scanner struct {
	buf     *bufio.Scanner
	ch      chan string
	started bool
	lastMsg string
	lastErr error
}

func New(r io.Reader) *Scanner {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(nil, bytesutil.MaxLineLength)
	return &Scanner{
		buf: scanner,
		ch:  make(chan string, 100),
	}
}

func (s *Scanner) scanGoroutine() {
	defer func() {
		close(s.ch)
		if err := recover(); err != nil {
			switch err := err.(type) {
			case error:
				s.lastErr = err
			default:
				s.lastErr = errors.New(fmt.Sprint(err))
			}
		}
	}()
	for s.buf.Scan() {
		s.ch <- s.buf.Text()
	}
	if err := s.buf.Err(); err != nil {
		s.lastErr = err
	}
}

func (s *Scanner) Scan(ctx context.Context) (bool, error) {
	if s.lastErr != nil {
		return false, s.lastErr
	}
	if err := ctx.Err(); err != nil {
		// Don't start the scanner if the context is already cancelled
		return false, err
	}
	if !s.started {
		s.started = true
		go s.scanGoroutine()
	}
	select {
	case text, ok := <-s.ch:
		// ok=false when channel is closed
		if !ok {
			return false, s.lastErr
		}
		s.lastMsg = text
		return true, nil
	case <-ctx.Done():
		return false, ctx.Err()
	}
}

func (s *Scanner) Text() string {
	return s.lastMsg
}

func (s *Scanner) Bytes() []byte {
	return []byte(s.lastMsg)
}

func (s *Scanner) Err() error {
	return s.lastErr
}
