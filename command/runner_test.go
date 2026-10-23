package command

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"

	"github.com/kubecolor/kubecolor/config"
)

func TestExec_kubectlNotFound(t *testing.T) {
	cfg := &Config{Config: &config.Config{
		Kubectl: "foo-bar-some-executable-that-does-not-exist",
	}}

	tests := []struct {
		name string
		run  func() error
	}{
		{name: "colored", run: func() error { return runColored(cfg) }},
		{name: "plain", run: func() error { return execWithoutColors(cfg, nil) }},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.run()
			if !errors.Is(err, exec.ErrNotFound) {
				t.Fatalf("want exec.ErrNotFound, got: %v", err)
			}
			if _, ok := errors.AsType[*KubectlError](err); ok {
				t.Errorf("want error other than *KubectlError, got: %v", err)
			}
			suffix := "; kubectl must be installed to use kubecolor"
			if !strings.HasSuffix(err.Error(), suffix) {
				t.Errorf("Wrong error\nwant suffix: %q\ngot: %q", suffix, err)
			}
		})
	}
}

func TestExec_kubectlKilledBySignal(t *testing.T) {
	cfg := &Config{Config: &config.Config{Kubectl: fakeKubectl(t, "kill -TERM $$")}}

	tests := []struct {
		name string
		run  func() error
	}{
		{name: "colored", run: func() error { return runColored(cfg) }},
		{name: "plain", run: func() error { return execWithoutColors(cfg, nil) }},
		{name: "complete", run: func() error {
			_, err := runKubectlComplete(nil, cfg)
			return err
		}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.run()
			ke, ok := errors.AsType[*KubectlError](err)
			if !ok {
				t.Fatalf("want *KubectlError, got: %v", err)
			}
			if ke.ExitCode != 143 {
				t.Errorf("want ExitCode 143, got: %d", ke.ExitCode)
			}
			exitErr, ok := errors.AsType[*exec.ExitError](ke)
			if !ok {
				t.Fatalf("want *KubectlError wrapping *exec.ExitError, got: %v", ke)
			}
			if ws := exitErr.Sys().(syscall.WaitStatus); ws.Signal() != syscall.SIGTERM {
				t.Errorf("want killed by SIGTERM, got: %v", exitErr)
			}
		})
	}
}

func TestExec_waitErrorNotKubectlError(t *testing.T) {
	tests := []struct {
		name    string
		run     func(t *testing.T) error
		wantErr string
	}{
		{
			name: "plain stdout write fails",
			run: func(t *testing.T) error {
				// kubectl exits 0, but its output can't be written. That must come back as a plain error: as a
				// KubectlError it'd carry exit code 0, so kubecolor would exit 0 without reporting anything. exec copies
				// output into a non-file Stdout (e.g. colorable's Windows console writer) itself, and Wait returns
				// that copy's error.
				cfg := &Config{Config: &config.Config{Kubectl: fakeKubectl(t, "echo apple")}}
				oldStdout := Stdout
				Stdout = failingWriter{errors.New("banana")}
				t.Cleanup(func() { Stdout = oldStdout })
				return execWithoutColors(cfg, nil)
			},
			wantErr: "banana",
		},
		{
			name: "colored never started",
			run: func(t *testing.T) error {
				// Close must also return non-exit Wait errors as is, rather than as a KubectlError with a made-up exit
				// code. Real ones come from the OS and can't be triggered here (the colored path's pipes are files, so
				// there's no copy to fail), so waiting on a never-started command stands in.
				r := &cmdWaitReadCloser{cmd: exec.Command("cherry"), stdout: io.NopCloser(strings.NewReader(""))}
				return r.Close()
			},
			wantErr: "exec: not started",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.run(t)
			// Exact match, so the error must come back with nothing wrapped around it.
			if err == nil || err.Error() != tc.wantErr {
				t.Fatalf("want error %q, got: %v", tc.wantErr, err)
			}
			if _, ok := errors.AsType[*KubectlError](err); ok {
				t.Errorf("want error other than *KubectlError, got: %v", err)
			}
		})
	}
}

// fakeKubectl writes a /bin/sh script to stand in for kubectl and returns its path.
//
// Skips the test on Windows, which can't run it.
func fakeKubectl(t *testing.T, script string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("requires /bin/sh")
	}
	path := filepath.Join(t.TempDir(), "kubectl")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+script+"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

type failingWriter struct{ err error }

func (w failingWriter) Write([]byte) (int, error) { return 0, w.err }

// runColored runs kubectl via the colored path and returns its error, as execWithoutColors does for the plain path.
//
// The colored path reports errors in two places: execWithReaders if kubectl fails to start, and closing stdout,
// which waits for kubectl and reports how it exited. Like Run, this reads stdout to the end before closing it, but
// discards the output. Stderr is never read, so the fake kubectl must not write more than a pipe buffer to it.
func runColored(cfg *Config) error {
	stdout, _, err := execWithReaders(cfg, nil)
	if err != nil {
		return err
	}
	if _, err := io.Copy(io.Discard, stdout); err != nil {
		return err
	}
	return stdout.Close()
}
