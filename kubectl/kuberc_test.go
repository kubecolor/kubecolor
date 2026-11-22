package kubectl

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

const testKuberc = `apiVersion: kubectl.config.k8s.io/v1beta1
kind: Preference
aliases:
- name: g
  command: get
- name: d
  command: describe
- name: gpa
  command: get
  prependArgs:
  - pods
  options:
  - name: all-namespaces
    default: "true"
  - name: output
    default: wide
- name: crns
  command: create namespace
- name: gj
  command: get
  appendArgs:
  - -o
  - json
- name: get
  command: describe
`

func TestExpandKubercAlias(t *testing.T) {
	kuberc := &Kuberc{}
	if err := yaml.Unmarshal([]byte(testKuberc), kuberc); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		args string
		want string
	}{
		{"g pods", "get pods"},
		{"d pod my-pod", "describe pod my-pod"},
		{"gpa", "get pods --all-namespaces=true --output=wide"},
		{"gpa -o yaml", "get pods --all-namespaces=true --output=wide -o yaml"},
		{"crns foo", "create namespace foo"},
		{"gj pods", "get pods -o json"},
		{"--context=dev g pods", "--context=dev get pods"},
		// built-in commands are never aliased
		{"get pods", "get pods"},
		// unknown names are left as-is
		{"foo bar", "foo bar"},
		{"", ""},
	}

	for _, tc := range tests {
		t.Run(tc.args, func(t *testing.T) {
			got := strings.Join(ExpandKubercAlias(strings.Fields(tc.args), kuberc), " ")
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestExpandKubercAlias_SubcommandInfo(t *testing.T) {
	kuberc := &Kuberc{}
	if err := yaml.Unmarshal([]byte(testKuberc), kuberc); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		args     string
		expected *SubcommandInfo
	}{
		{"g pods", &SubcommandInfo{Subcommand: Get}},
		{"gpa", &SubcommandInfo{Subcommand: Get, Output: OutputWide}},
		{"gpa -o json", &SubcommandInfo{Subcommand: Get, Output: OutputJSON}},
		{"gj pods", &SubcommandInfo{Subcommand: Get, Output: OutputJSON}},
		{"d pod my-pod", &SubcommandInfo{Subcommand: Describe}},
	}

	for _, tc := range tests {
		t.Run(tc.args, func(t *testing.T) {
			args := ExpandKubercAlias(strings.Fields(tc.args), kuberc)
			got := InspectSubcommandInfo(args, NoopPluginHandler{})
			if diff := cmp.Diff(tc.expected, got); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}

func TestLoadKuberc(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "kuberc")
	if err := os.WriteFile(path, []byte(testKuberc), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", dir)
	t.Setenv("USERPROFILE", dir)
	t.Setenv("KUBECTL_KUBERC", "")

	t.Run("flag", func(t *testing.T) {
		t.Setenv("KUBERC", "")
		if k := LoadKuberc([]string{"--kuberc", path, "g"}); k == nil || len(k.Aliases) == 0 {
			t.Errorf("--kuberc %s: no aliases loaded", path)
		}
		if k := LoadKuberc([]string{"--kuberc=" + path, "g"}); k == nil || len(k.Aliases) == 0 {
			t.Errorf("--kuberc=%s: no aliases loaded", path)
		}
	})

	t.Run("env", func(t *testing.T) {
		t.Setenv("KUBERC", path)
		if k := LoadKuberc([]string{"g"}); k == nil || len(k.Aliases) == 0 {
			t.Error("KUBERC: no aliases loaded")
		}
	})

	t.Run("disabled", func(t *testing.T) {
		t.Setenv("KUBERC", "off")
		if k := LoadKuberc([]string{"g"}); k != nil {
			t.Error("KUBERC=off: kuberc was loaded")
		}
		t.Setenv("KUBERC", path)
		t.Setenv("KUBECTL_KUBERC", "false")
		if k := LoadKuberc([]string{"g"}); k != nil {
			t.Error("KUBECTL_KUBERC=false: kuberc was loaded")
		}
	})

	t.Run("default", func(t *testing.T) {
		t.Setenv("KUBERC", "")
		if k := LoadKuberc([]string{"g"}); k != nil {
			t.Error("no ~/.kube/kuberc: kuberc was loaded")
		}
		if err := os.MkdirAll(filepath.Join(dir, ".kube"), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, ".kube", "kuberc"), []byte(testKuberc), 0o600); err != nil {
			t.Fatal(err)
		}
		if k := LoadKuberc([]string{"g"}); k == nil || len(k.Aliases) == 0 {
			t.Error("~/.kube/kuberc: no aliases loaded")
		}
	})

	t.Run("missing", func(t *testing.T) {
		t.Setenv("KUBERC", filepath.Join(dir, "does-not-exist"))
		if k := LoadKuberc([]string{"g"}); k != nil {
			t.Error("missing file: kuberc was loaded")
		}
	})
}
