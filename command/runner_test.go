package command

import (
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/config"
)

func TestExecWithReaders_notFound(t *testing.T) {
	r, w, err := execWithReaders(&Config{Config: &config.Config{
		Kubectl: "foo-bar-some-executable-that-does-not-exist",
	}}, []string{})
	if err == nil {
		defer r.Close()
		defer w.Close()
		t.Error("Expected error, but got nil")
		return
	}

	suffix := "; kubectl must be installed to use kubecolor"
	if !strings.HasSuffix(err.Error(), suffix) {
		t.Errorf("Wrong error\nwant suffix: %q\ngot: %q", suffix, err)
	}
}
