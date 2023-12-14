package describe

import (
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/testutil"
)

func TestScanner_noKeyIndent(t *testing.T) {
	const input = "" +
		"Name:             traefik-64d54f8757-blrj9\n" +
		"Namespace:        traefik\n" +
		"Labels:           app.kubernetes.io/instance=traefik-traefik\n" +
		"                  app.kubernetes.io/managed-by=Helm\n" +
		"                  app.kubernetes.io/name=traefik\n" +
		"Status:           Running\n" +
		"IP:               10.0.0.1\n"

	s := New(strings.NewReader(input))

	mustScanToken(t, s, KindKey, "Name:", 0, 18)
	mustScanToken(t, s, KindWhitespace, "             ", 0, 18)
	mustScanToken(t, s, KindValue, "traefik-64d54f8757-blrj9", 0, 18)
	mustScanToken(t, s, KindEOL, "\n", 0, 18)

	mustScanToken(t, s, KindKey, "Namespace:", 0, 18)
	mustScanToken(t, s, KindWhitespace, "        ", 0, 18)
	mustScanToken(t, s, KindValue, "traefik", 0, 18)
	mustScanToken(t, s, KindEOL, "\n", 0, 18)

	mustScanToken(t, s, KindKey, "Labels:", 0, 18)
	mustScanToken(t, s, KindWhitespace, "           ", 0, 18)
	mustScanToken(t, s, KindValue, "app.kubernetes.io/instance=traefik-traefik", 0, 18)
	mustScanToken(t, s, KindEOL, "\n", 0, 18)

	mustScanToken(t, s, KindWhitespace, "                  ", 0, 18)
	mustScanToken(t, s, KindValue, "app.kubernetes.io/managed-by=Helm", 0, 18)
	mustScanToken(t, s, KindEOL, "\n", 0, 18)

	mustScanToken(t, s, KindWhitespace, "                  ", 0, 18)
	mustScanToken(t, s, KindValue, "app.kubernetes.io/name=traefik", 0, 18)
	mustScanToken(t, s, KindEOL, "\n", 0, 18)

	mustScanToken(t, s, KindKey, "Status:", 0, 18)
	mustScanToken(t, s, KindWhitespace, "           ", 0, 18)
	mustScanToken(t, s, KindValue, "Running", 0, 18)
	mustScanToken(t, s, KindEOL, "\n", 0, 18)

	mustScanToken(t, s, KindKey, "IP:", 0, 18)
	mustScanToken(t, s, KindWhitespace, "               ", 0, 18)
	mustScanToken(t, s, KindValue, "10.0.0.1", 0, 18)
	mustScanToken(t, s, KindEOL, "\n", 0, 18)

	if s.Scan() {
		t.Fatalf("Expected no more scans, but got: %#v", s.Token())
	}
}

func mustScanToken(t *testing.T, s *Scanner, wantKind Kind, wantString string, wantKeyIndent int, wantValueIndent int) {
	t.Helper()
	if !s.Scan() {
		t.Fatalf("Failed scan; Expected value kind=%s: %q", wantKind, wantString)
	}
	testutil.MustEqual(t, Token{
		Kind:  wantKind,
		Bytes: []byte(wantString),
	}, s.Token())
	if wantKeyIndent != s.KeyIndent() {
		t.Fatalf("Wrong value; Want key indent %d, got %d", wantKeyIndent, s.KeyIndent())
	}
	if wantValueIndent != s.ValueIndent() {
		t.Fatalf("Wrong value; Want key indent %d, got %d", wantValueIndent, s.ValueIndent())
	}
}
