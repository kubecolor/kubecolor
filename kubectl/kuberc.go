package kubectl

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Kuberc is the subset of kubectl's kuberc file (kubectl.config.k8s.io
// Preference) that kubecolor needs to know which command an alias runs.
type Kuberc struct {
	Aliases []KubercAlias `yaml:"aliases"`
}

type KubercAlias struct {
	Name        string                `yaml:"name"`
	Command     string                `yaml:"command"`
	PrependArgs []string              `yaml:"prependArgs"`
	AppendArgs  []string              `yaml:"appendArgs"`
	Options     []KubercOptionDefault `yaml:"options"`
}

type KubercOptionDefault struct {
	Name    string `yaml:"name"`
	Default string `yaml:"default"`
}

// LoadKuberc reads the kuberc file kubectl would use for the given args.
// It returns nil when kuberc is disabled, missing or invalid, in which case
// kubectl reports the problem itself.
func LoadKuberc(args []string) *Kuberc {
	path := kubercPath(args)
	if path == "" {
		return nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var kuberc Kuberc
	if err := yaml.Unmarshal(b, &kuberc); err != nil {
		return nil
	}
	return &kuberc
}

// kubercPath follows the lookup order of kubectl (k8s.io/kubectl/pkg/kuberc):
// the --kuberc flag, then $KUBERC, then ~/.kube/kuberc.
// KUBERC=off and KUBECTL_KUBERC=false disable it.
func kubercPath(args []string) string {
	if os.Getenv("KUBERC") == "off" || os.Getenv("KUBECTL_KUBERC") == "false" {
		return ""
	}
	for i, arg := range args {
		if arg == "--" {
			break
		}
		if arg == "--kuberc" {
			if i+1 < len(args) {
				return args[i+1]
			}
			return ""
		}
		if path, ok := strings.CutPrefix(arg, "--kuberc="); ok {
			return path
		}
	}
	if path := os.Getenv("KUBERC"); path != "" {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".kube", "kuberc")
}

// ExpandKubercAlias replaces a kuberc alias in args with the command it
// stands for, the same way kubectl does, so the real subcommand and flags
// can be inspected. Args are returned unchanged when no alias applies.
func ExpandKubercAlias(args []string, kuberc *Kuberc) []string {
	if kuberc == nil {
		return args
	}

	// kubectl takes the first argument that is not a flag as the alias name
	index := -1
	for i, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			index = i
			break
		}
	}
	if index == -1 {
		return args
	}
	name := args[index]

	// aliases can not shadow built-in commands
	if _, ok := InspectSubcommand([]string{name}, NoopPluginHandler{}); ok {
		return args
	}

	for _, alias := range kuberc.Aliases {
		if alias.Name != name {
			continue
		}
		command := strings.Fields(alias.Command)
		if len(command) == 0 {
			return args
		}
		expanded := make([]string, 0, len(args)+len(command)+len(alias.PrependArgs)+len(alias.Options)+len(alias.AppendArgs))
		expanded = append(expanded, args[:index]...)
		expanded = append(expanded, command...)
		expanded = append(expanded, alias.PrependArgs...)
		// alias options are defaults, so they go before the user's own
		// args, which then take precedence when parsed
		for _, opt := range alias.Options {
			expanded = append(expanded, "--"+opt.Name+"="+opt.Default)
		}
		expanded = append(expanded, args[index+1:]...)
		expanded = append(expanded, alias.AppendArgs...)
		return expanded
	}
	return args
}
