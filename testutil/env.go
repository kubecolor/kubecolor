package testutil

import (
	"os"
	"testing"
)

// Setenv is a helper function for setting environment variables in tests.
// If the value is empty, then the environment variable is unset.
// At the end of the test (using [testing.TB.Cleanup]) the environment variable
// is unset.
func Setenv(t testing.TB, key, value string) {
	t.Helper()
	if value == "" {
		os.Unsetenv(key)
		return
	}

	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("Failed to set environment variable %q: %s", key, err)
	}

	t.Cleanup(func() {
		os.Unsetenv(key)
	})
}
