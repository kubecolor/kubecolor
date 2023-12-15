package describe

import (
	"bytes"
	"embed"
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/testutil"
)

//go:embed testdata
var testdataFS embed.FS

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

	s := NewScanner(strings.NewReader(input))

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

	s := NewScanner(strings.NewReader(input))

	mustScanLine(t, s, "~Args:~~~", "Args")
	mustScanLine(t, s, "  ~~~--first-flag~", "Args")
	mustScanLine(t, s, "  ~~~--second-flag=with single spaces in value~", "Args")
	mustScanLine(t, s, "  ~~~--last-flag~", "Args")

	if s.Scan() {
		t.Fatalf("Expected no more scans, but got: %#v", s.Line())
	}
}

func TestScanner_nested(t *testing.T) {
	const input = "" +
		"Containers:\n" +
		"  traefik:\n" +
		"    Image:          docker.io/traefik:v2.10.6\n" +
		"    State:          Running\n" +
		"      Started:      Wed, 13 Dec 2023 22:55:39 +0100\n" +
		"  linkerd:\n" +
		"    State:          Running\n"

	s := NewScanner(strings.NewReader(input))

	mustScanLine(t, s, "~Containers:~~~", "Containers")
	mustScanLine(t, s, "  ~traefik:~~~", "Containers/traefik")
	mustScanLine(t, s, "    ~Image:~          ~docker.io/traefik:v2.10.6~", "Containers/traefik/Image")
	mustScanLine(t, s, "    ~State:~          ~Running~", "Containers/traefik/State")
	mustScanLine(t, s, "      ~Started:~      ~Wed, 13 Dec 2023 22:55:39 +0100~", "Containers/traefik/State/Started")
	mustScanLine(t, s, "  ~linkerd:~~~", "Containers/linkerd")
	mustScanLine(t, s, "    ~State:~          ~Running~", "Containers/linkerd/State")

	if s.Scan() {
		t.Fatalf("Expected no more scans, but got: %#v", s.Line())
	}
}

func TestScanner_tabbedValues(t *testing.T) {
	const input = "" +
		"  apiVersion\t<string>\n" +
		"  kind\t<string>\n" +
		"  metadata\t<ObjectMeta>\n" +
		"    name\t<string>\n" +
		"    labels\t<map[string]string>\n" +
		"  spec\t<PodSpec>\n" +
		"    containers\t<[]Container> -required-\n" +
		"      args\t<[]string>\n" +
		"      command\t<[]string>\n" +
		"      env\t<[]EnvVar>\n" +
		"        name\t<string> -required-\n" +
		"        value\t<string>\n"

	s := NewScanner(strings.NewReader(input))

	mustScanLine(t, s, "  ~apiVersion~\t~<string>~", "apiVersion")
	mustScanLine(t, s, "  ~kind~\t~<string>~", "kind")
	mustScanLine(t, s, "  ~metadata~\t~<ObjectMeta>~", "metadata")
	mustScanLine(t, s, "    ~name~\t~<string>~", "metadata/name")
	mustScanLine(t, s, "    ~labels~\t~<map[string]string>~", "metadata/labels")
	mustScanLine(t, s, "  ~spec~\t~<PodSpec>~", "spec")
	mustScanLine(t, s, "    ~containers~\t~<[]Container> -required-~", "spec/containers")
	mustScanLine(t, s, "      ~args~\t~<[]string>~", "spec/containers/args")
	mustScanLine(t, s, "      ~command~\t~<[]string>~", "spec/containers/command")
	mustScanLine(t, s, "      ~env~\t~<[]EnvVar>~", "spec/containers/env")
	mustScanLine(t, s, "        ~name~\t~<string> -required-~", "spec/containers/env/name")
	mustScanLine(t, s, "        ~value~\t~<string>~", "spec/containers/env/value")

	if s.Scan() {
		t.Fatalf("Expected no more scans, but got: %#v", s.Line())
	}
}

func TestScanner_explainText(t *testing.T) {
	const input = "" +
		"FIELDS:\n" +
		"  apiVersion\t<string>\n" +
		"    APIVersion defines the versioned schema of this representation of an object.\n" +
		"    Servers should convert recognized schemas to the latest internal value, and\n" +
		"    may reject unrecognized values. More info:\n" +
		"    https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources\n" +
		"  spec\t<PodSpec>\n" +
		"    Specification of the desired behavior of the pod. More info:\n" +
		"    https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status\n"

	s := NewScanner(strings.NewReader(input))

	mustScanLine(t, s, "~FIELDS:~~~", "FIELDS")
	mustScanLine(t, s, "  ~apiVersion~\t~<string>~", "FIELDS/apiVersion")
	mustScanLine(t, s, "    ~~~APIVersion defines the versioned schema of this representation of an object.~", "FIELDS/apiVersion")
	mustScanLine(t, s, "    ~~~Servers should convert recognized schemas to the latest internal value, and~", "FIELDS/apiVersion")
	mustScanLine(t, s, "    ~~~may reject unrecognized values. More info:~", "FIELDS/apiVersion")
	mustScanLine(t, s, "    ~~~https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources~", "FIELDS/apiVersion")
	mustScanLine(t, s, "  ~spec~\t~<PodSpec>~", "FIELDS/spec")
	mustScanLine(t, s, "    ~~~Specification of the desired behavior of the pod. More info:~", "FIELDS/spec")
	mustScanLine(t, s, "    ~~~https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status~", "FIELDS/spec")

	if s.Scan() {
		t.Fatalf("Expected no more scans, but got: %#v", s.Line())
	}
}

func TestScanner_recreateString(t *testing.T) {
	b, err := testdataFS.ReadFile("testdata/explain-pod.txt")
	if err != nil {
		t.Fatalf("Read testdata file: %s", err)
	}

	var buf bytes.Buffer

	s := NewScanner(bytes.NewReader(b))
	for s.Scan() {
		buf.WriteString(s.Line().String())
		buf.WriteByte('\n')
	}
	if err := s.Err(); err != nil {
		t.Fatalf("Scan error: %s", err)
	}

	if buf.Len() < 39_000 {
		t.Fatalf("The file content should be about 40,000 bytes, but got %d", buf.Len())
	}

	testutil.MustEqual(t, string(b), buf.String())
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

	gotPath := s.Path().String()
	if gotPath != path {
		t.Errorf("Wrong path; Expected %q, but got %q", path, gotPath)
	}
}
