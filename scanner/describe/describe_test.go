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
	mustMatchState(t, s, 0, 18, []string{"Name"})

	mustScanToken(t, s, KindKey, "Namespace:", 0, 18)
	mustScanToken(t, s, KindWhitespace, "        ", 0, 18)
	mustScanToken(t, s, KindValue, "traefik", 0, 18)
	mustScanToken(t, s, KindEOL, "\n", 0, 18)
	mustMatchState(t, s, 0, 18, []string{"Namespace"})

	mustScanToken(t, s, KindKey, "Labels:", 0, 18)
	mustScanToken(t, s, KindWhitespace, "           ", 0, 18)
	mustScanToken(t, s, KindValue, "app.kubernetes.io/instance=traefik-traefik", 0, 18)
	mustScanToken(t, s, KindEOL, "\n", 0, 18)
	mustMatchState(t, s, 0, 18, []string{"Labels"})

	mustScanToken(t, s, KindWhitespace, "                  ", 0, 18)
	mustScanToken(t, s, KindValue, "app.kubernetes.io/managed-by=Helm", 0, 18)
	mustScanToken(t, s, KindEOL, "\n", 0, 18)
	mustMatchState(t, s, 0, 18, []string{"Labels"})

	mustScanToken(t, s, KindWhitespace, "                  ", 0, 18)
	mustScanToken(t, s, KindValue, "app.kubernetes.io/name=traefik", 0, 18)
	mustScanToken(t, s, KindEOL, "\n", 0, 18)
	mustMatchState(t, s, 0, 18, []string{"Labels"})

	mustScanToken(t, s, KindKey, "Status:", 0, 18)
	mustScanToken(t, s, KindWhitespace, "           ", 0, 18)
	mustScanToken(t, s, KindValue, "Running", 0, 18)
	mustScanToken(t, s, KindEOL, "\n", 0, 18)
	mustMatchState(t, s, 0, 18, []string{"Status"})

	mustScanToken(t, s, KindKey, "IP:", 0, 18)
	mustScanToken(t, s, KindWhitespace, "               ", 0, 18)
	mustScanToken(t, s, KindValue, "10.0.0.1", 0, 18)
	mustScanToken(t, s, KindEOL, "\n", 0, 18)
	mustMatchState(t, s, 0, 18, []string{"IP"})

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
}

func mustMatchState(t *testing.T, s *Scanner, wantKeyIndent, wantValueIndent int, path []string) {
	testutil.MustEqual(t, State{
		KeyIndent:   wantKeyIndent,
		ValueIndent: wantValueIndent,
		Path:        path,
	}, s.State())
}
