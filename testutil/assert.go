// testutil package is a utility for testing.
// This package is inspired by morikuni/failure testutil_test.go
package testutil

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Equal(t testing.TB, want, got any, msg ...any) {
	t.Helper()

	if diff := cmp.Diff(want, got); diff != "" {
		if len(msg) > 0 {
			t.Errorf("%s\ndiff (-want +got):\n%s", fmt.Sprint(msg...), diff)
		} else {
			t.Errorf("diff (-want +got):\n%s", diff)
		}
	}
}

func Equalf(t testing.TB, want, got any, format string, msg ...any) {
	t.Helper()

	if diff := cmp.Diff(want, got); diff != "" {
		if format != "" {
			t.Errorf("%s\ndiff (-want +got):\n%s", fmt.Sprintf(format, msg...), diff)
		} else {
			t.Errorf("diff (-want +got):\n%s", diff)
		}
	}
}

func MustEqual(t testing.TB, want, got any, msg ...any) {
	t.Helper()

	if diff := cmp.Diff(want, got); diff != "" {
		if len(msg) > 0 {
			t.Fatalf("%s\ndiff (-want +got):\n%s", fmt.Sprint(msg...), diff)
		} else {
			t.Fatalf("diff (-want +got):\n%s", diff)
		}
	}
}

func MustEqualf(t testing.TB, want, got any, format string, msg ...any) {
	t.Helper()

	if diff := cmp.Diff(want, got); diff != "" {
		if format != "" {
			t.Fatalf("%s\ndiff (-want +got):\n%s", fmt.Sprintf(format, msg...), diff)
		} else {
			t.Fatalf("diff (-want +got):\n%s", diff)
		}
	}
}

func MustNoError(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func NoError(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}
