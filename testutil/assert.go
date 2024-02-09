// testutil package is a utility for testing.
// This package is inspired by morikuni/failure testutil_test.go
package testutil

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func MustEqual(t testing.TB, want, got any) {
	t.Helper()

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("diff (-want +got):\n%s", diff)
	}
}

func NoError(t testing.TB, err error, unexpectedValue any) {
	t.Helper()

	if err != nil {
		t.Errorf("unexpected error: %s\nunexpected value: %#v", err, unexpectedValue)
	}
}
