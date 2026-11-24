// Copyright 2025 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// This file contains a modified version of kubectl's kuberc lookup, from the
// LoadKuberc and getExplicitKuberc functions:
// [https://github.com/kubernetes/kubectl/blob/b6fd7151e781c4619468e46cfefb5f4728cce747/pkg/kuberc/kuberc.go#L357-L377]
// [https://github.com/kubernetes/kubectl/blob/b6fd7151e781c4619468e46cfefb5f4728cce747/pkg/kuberc/kuberc.go#L424-L453]
//
// As well as of the alias expansion in Preferences.applyAliases:
// [https://github.com/kubernetes/kubectl/blob/b6fd7151e781c4619468e46cfefb5f4728cce747/pkg/kuberc/kuberc.go#L230-L352]

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
	Name        string         `yaml:"name"`
	Command     string         `yaml:"command"`
	PrependArgs []string       `yaml:"prependArgs"`
	AppendArgs  []string       `yaml:"appendArgs"`
	Options     []KubercOption `yaml:"options"`
}

type KubercOption struct {
	Name    string `yaml:"name"`
	Default string `yaml:"default"`
}

// LoadKuberc reads the kuberc file kubectl would use for the given args.
// It returns an empty Kuberc when kuberc is disabled, missing or invalid,
// in which case kubectl reports the problem itself.
func LoadKuberc(args []string) Kuberc {
	path := kubercPath(args)
	if path == "" {
		return Kuberc{}
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return Kuberc{}
	}
	var kuberc Kuberc
	if err := yaml.Unmarshal(b, &kuberc); err != nil {
		return Kuberc{}
	}
	return kuberc
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

// ExpandAlias replaces a kuberc alias in args with the command it stands
// for, the same way kubectl does, so the real subcommand and flags can be
// inspected. Args are returned unchanged when no alias applies.
func (k Kuberc) ExpandAlias(args []string) []string {
	index, ok := firstNonFlagArg(args)
	if !ok {
		return args
	}
	name := args[index]

	// aliases can not shadow built-in commands
	if _, ok := InspectSubcommand([]string{name}, NoopPluginHandler{}); ok {
		return args
	}

	alias, ok := k.FindAlias(name)
	if !ok {
		return args
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

// FindAlias returns the alias with the given name.
func (k Kuberc) FindAlias(name string) (KubercAlias, bool) {
	for _, alias := range k.Aliases {
		if alias.Name == name {
			return alias, true
		}
	}
	return KubercAlias{}, false
}

// firstNonFlagArg returns the index of the first argument that is not a
// flag, which is the one kubectl takes as the alias name.
func firstNonFlagArg(args []string) (int, bool) {
	for i, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			return i, true
		}
	}
	return -1, false
}
