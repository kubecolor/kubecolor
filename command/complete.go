package command

import (
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/kubecolor/kubecolor/kubectl"
)

type InjectCompletionsOptions struct {
	RawArgs []string
	Args    []string
	Config  *Config
}

func InjectKubecolorCompletions(rawArgs []string, cfg *Config, sci *kubectl.SubcommandInfo) error {
	toComplete := rawArgs[len(rawArgs)-1]

	// We can't feed kubecolor flags to kubectl, as then it would error.
	// So we hack around it by using the args slice without the kubecolor flags.
	// But then we must translate `kubectl __complete --plain`
	// so it's sent to kubectl as `kubectl __complete -`
	kubectlArgs := cfg.ArgsPassthrough
	if len(kubectlArgs) == 1 && strings.HasPrefix(toComplete, "-") {
		kubectlArgs = append(kubectlArgs, "-")
	}

	s, err := runKubectlComplete(kubectlArgs, cfg)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(toComplete, "-") {
		fmt.Print(s)
		return nil
	}

	output := ParseKubectlComplete(s)
	args, directive := GetFlagCompletions(cfg.Flags, toComplete)
	output.Args = append(output.Args, args...)
	if directive != 0 {
		output.Directive = directive
	}

	// We have to do additional filtering here because of our hack above
	output.Args = slices.DeleteFunc(output.Args, func(arg CompleteArg) bool {
		return !strings.HasPrefix(arg.Name, toComplete)
	})

	slices.SortFunc(args, func(a, b CompleteArg) int {
		return cmp.Compare(a.Name, b.Name)
	})

	noDesc := sci.Subcommand == kubectl.CompleteNoDesc
	for _, arg := range output.Args {
		if noDesc {
			fmt.Println(arg.Name)
		} else {
			fmt.Printf("%s\t%s\n", arg.Name, arg.Description)
		}
	}

	fmt.Printf(":%d\n", output.Directive)

	return nil
}

func GetFlagCompletions(flags FlagSet, toComplete string) ([]CompleteArg, CompleteDirective) {
	directive := CompleteDirectiveDefault
	args := make([]CompleteArg, 0, len(flags))
	addIfSingle := ""

	for _, flag := range flags {
		if !strings.HasPrefix(flag.Name, toComplete) {
			continue
		}
		name := flag.Name
		if flag.RequiresValue {
			directive = CompleteDirectiveNoSpace
			addIfSingle = "="
		}
		args = append(args, CompleteArg{
			Name:        name,
			Description: flag.Description,
		})
	}

	if len(args) == 1 && addIfSingle != "" {
		args[0].Name += addIfSingle
	}

	return args, directive
}

type CompleteDirective int

const (
	CompleteDirectiveError CompleteDirective = 1 << iota
	CompleteDirectiveNoSpace
	CompleteDirectiveNoFileComp
	CompleteDirectiveFilterFileExt
	CompleteDirectiveFilterDirs
	CompleteDirectiveKeepOrder

	CompleteDirectiveDefault CompleteDirective = 0
)

type CompleteOutput struct {
	Args      []CompleteArg
	Directive CompleteDirective
}

type CompleteArg struct {
	Name        string
	Description string
}

func ParseKubectlComplete(s string) CompleteOutput {
	var output CompleteOutput
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, ":") {
			var code int
			fmt.Sscanf(line, ":%d", &code)
			output.Directive = CompleteDirective(code)
			break
		}

		name, desc, _ := strings.Cut(line, "\t")
		output.Args = append(output.Args, CompleteArg{name, desc})
	}
	return output
}

func runKubectlComplete(rawArgs []string, cfg *Config) (string, error) {
	cmd := exec.Command(cfg.Kubectl, rawArgs...)

	var stdout bytes.Buffer
	cmd.Stdin = os.Stdin
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		os.Stdout.Write(stdout.Bytes())
		var execErr *exec.ExitError
		if errors.As(err, &execErr) {
			return "", &KubectlError{ExitCode: execErr.ExitCode()}
		}
		return "", err
	}

	return stdout.String(), nil
}
