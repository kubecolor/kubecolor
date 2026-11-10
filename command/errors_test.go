package command

import (
	"errors"
	"os/exec"
	"runtime"
	"testing"
)

func TestKubectlError(t *testing.T) {
	tests := []struct {
		name string
		err  *KubectlError
		want string
	}{
		{
			name: "exit code only",
			err:  &KubectlError{ExitCode: 1},
			want: "kubectl error: 1",
		},
		{
			name: "with cause",
			err:  &KubectlError{ExitCode: 143, Err: errors.New("signal: terminated")},
			want: "kubectl error: signal: terminated",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.err.Error(); got != tc.want {
				t.Errorf("want %q, got %q", tc.want, got)
			}
		})
	}
}

func TestWrapExitError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("requires POSIX signals and /bin/sh")
	}

	tests := []struct {
		name   string
		script string
		want   int
	}{
		{name: "exit status", script: "exit 3", want: 3},
		{name: "SIGTERM", script: "kill -TERM $$", want: 143},
		{name: "SIGKILL", script: "kill -KILL $$", want: 137},
		{name: "SIGPIPE", script: "kill -PIPE $$", want: 141},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := wrapExitError(exec.Command("/bin/sh", "-c", tc.script).Run())
			ke, ok := errors.AsType[*KubectlError](err)
			if !ok {
				t.Fatalf("want *KubectlError, got: %v", err)
			}
			if ke.ExitCode != tc.want {
				t.Errorf("want ExitCode %d, got %d", tc.want, ke.ExitCode)
			}
		})
	}
}
