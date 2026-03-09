package testutil

import "io"

type DummyReader struct {
	ReadFunc func([]byte) (int, error)
}

var _ io.Reader = DummyReader{}

// Read implements [io.Reader].
func (d DummyReader) Read(p []byte) (n int, err error) {
	return d.ReadFunc(p)
}
