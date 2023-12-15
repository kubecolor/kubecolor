package describe

import (
	"strings"
	"testing"
)

func TestScanner_noKeyIndent(t *testing.T) {
	const input = "" +
		"    \n" +
		"Name:             traefik-64d54f8757-blrj9\n" +
		"Namespace:        traefik\n" +
		"Labels:           app.kubernetes.io/instance=traefik-traefik\n" +
		"                  app.kubernetes.io/managed-by=Helm\n" +
		"                  app.kubernetes.io/name=traefik\n" +
		"Status:           Running\n" +
		"IP:               10.0.0.1\n"

	s := New(strings.NewReader(input))

	mustScanLine(t, s, "~~~~    ", "")
	mustScanLine(t, s, "~Name:~             ~traefik-64d54f8757-blrj9~", "Name")
	mustScanLine(t, s, "~Namespace:~        ~traefik~", "Namespace")
	mustScanLine(t, s, "~Labels:~           ~app.kubernetes.io/instance=traefik-traefik~", "Labels")
	mustScanLine(t, s, "                  ~~~app.kubernetes.io/managed-by=Helm~", "Labels")
	mustScanLine(t, s, "                  ~~~app.kubernetes.io/name=traefik~", "Labels")
	mustScanLine(t, s, "~Status:~           ~Running~", "Status")
	mustScanLine(t, s, "~IP:~               ~10.0.0.1~", "IP")

	if s.Scan() {
		t.Fatalf("Expected no more scans, but got: %#v", s.Line())
	}
}

func TestScanner_list(t *testing.T) {
	const input = "" +
		"Args:\n" +
		"  --first-flag\n" +
		"  --second-flag=with single spaces in value\n" +
		"  --last-flag\n"

	s := New(strings.NewReader(input))

	mustScanLine(t, s, "~Args:~~~", "Args")
	mustScanLine(t, s, "  ~~~--first-flag~", "Args")
	mustScanLine(t, s, "  ~~~--second-flag=with single spaces in value~", "Args")
	mustScanLine(t, s, "  ~~~--last-flag~", "Args")

	if s.Scan() {
		t.Fatalf("Expected no more scans, but got: %#v", s.Line())
	}
}

func mustScanLine(t *testing.T, s *Scanner, line, path string) {
	t.Helper()
	if !s.Scan() {
		t.Fatalf("Failed scan; Expected value %q", line)
	}

	gotLine := s.Line().GoString()
	if gotLine != line {
		t.Errorf("Wrong line (format: indent~key~spacing~value~trailing)\nWant %q\nGot  %q", line, gotLine)
	}

	gotPath := s.Path()
	if gotPath != path {
		t.Errorf("Wrong path; Expected %q, but got %q", path, gotPath)
	}
}
