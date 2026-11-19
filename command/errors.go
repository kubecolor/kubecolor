package command

import (
	"errors"
	"fmt"
	"os/exec"
	"syscall"
)

// KubectlError reports that kubectl ran but exited unsuccessfully.
//
// ExitCode is what a shell would report: kubectl's exit status, or 128+n if signal n killed it. Err is the
// *exec.ExitError it was built from. Any other failure, such as kubectl not being installed, reaches callers as its
// original error rather than a KubectlError.
type KubectlError struct {
	ExitCode int
	Err      error
}

func (ke *KubectlError) Error() string {
	// Err is only nil if a caller (incorrectly) builds a KubectlError by hand
	if ke.Err == nil {
		return fmt.Sprintf("kubectl error: %d", ke.ExitCode)
	}
	return fmt.Sprintf("kubectl error: %v", ke.Err)
}

func (ke *KubectlError) Unwrap() error {
	return ke.Err
}

// wrapExitError wraps err in a KubectlError if it's an *exec.ExitError, and otherwise returns it as is.
func wrapExitError(err error) error {
	exitErr, ok := errors.AsType[*exec.ExitError](err)
	if !ok {
		return err
	}
	// If a signal killed kubectl, exit with 128+n like a shell does (e.g. 141 for SIGPIPE, 143 for SIGTERM),
	// so kubecolor's status matches a bare kubectl's. ExitCode alone gives -1 here, which os.Exit turns into 255.
	// See https://www.gnu.org/software/bash/manual/html_node/Exit-Status.html.
	code := exitErr.ExitCode()
	if ws, ok := exitErr.Sys().(syscall.WaitStatus); ok && ws.Signaled() {
		code = 128 + int(ws.Signal())
	}
	return &KubectlError{ExitCode: code, Err: err}
}
