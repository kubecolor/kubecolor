package kubectl

import (
	"strings"
)

type SubcommandInfo struct {
	Subcommand     Subcommand
	SubcommandArgs []string // args after the subcommand
	Output         Output   // flag: -o, --output
	NoHeader       bool     // flag: --no-header
	Watch          bool     // flag: -w, --watch
	Follow         bool     // flag: -f, --follow
	Help           bool     // flag: -h, --help
	Recursive      bool     // flag: --recursive
	Client         bool     // flag: --client
	List           bool     // flag: --list
}

// Output is an enum of different "--output=..." types.
type Output byte

const (
	OutputNone              Output = iota
	OutputWide                     // -o wide
	OutputJSON                     // -o json
	OutputYAML                     // -o yaml
	OutputCustomColumns            // -o custom-columns=...
	OutputCustomColumnsFile        // -o custom-columns-file=...
	OutputOther                    // e.g -o jsonpath=...
)

// ParseOutput parses a "--output ..." flag's type
func ParseOutput(value string) Output {
	// only consider value before "="
	// e.g "--output custom-columns=..."
	output, _, _ := strings.Cut(value, "=")
	switch output {
	case "wide":
		return OutputWide
	case "json":
		return OutputJSON
	case "yaml":
		return OutputYAML
	case "custom-columns":
		return OutputCustomColumns
	case "custom-columns-file":
		return OutputCustomColumnsFile
	default:
		return OutputOther
	}
}

type Subcommand string

const (
	APIResources   Subcommand = "api-resources"
	APIVersions    Subcommand = "api-versions"
	Annotate       Subcommand = "annotate"
	Apply          Subcommand = "apply"
	Attach         Subcommand = "attach"
	Auth           Subcommand = "auth"
	Autoscale      Subcommand = "autoscale"
	Certificate    Subcommand = "certificate"
	ClusterInfo    Subcommand = "cluster-info"
	Complete       Subcommand = "__complete"
	CompleteNoDesc Subcommand = "__completeNoDesc"
	Completion     Subcommand = "completion"
	Config         Subcommand = "config"
	Convert        Subcommand = "convert"
	Cordon         Subcommand = "cordon"
	Cp             Subcommand = "cp"
	Create         Subcommand = "create"
	Debug          Subcommand = "debug"
	Delete         Subcommand = "delete"
	Describe       Subcommand = "describe"
	Diff           Subcommand = "diff"
	Drain          Subcommand = "drain"
	Edit           Subcommand = "edit"
	Events         Subcommand = "events"
	Exec           Subcommand = "exec"
	Explain        Subcommand = "explain"
	Expose         Subcommand = "expose"
	Get            Subcommand = "get"
	KubectlPlugin  Subcommand = "(plugin)"
	Kustomize      Subcommand = "kustomize"
	Label          Subcommand = "label"
	Logs           Subcommand = "logs"
	Options        Subcommand = "options"
	Patch          Subcommand = "patch"
	Plugin         Subcommand = "plugin"
	PortForward    Subcommand = "port-forward"
	Proxy          Subcommand = "proxy"
	Replace        Subcommand = "replace"
	Rollout        Subcommand = "rollout"
	Run            Subcommand = "run"
	Scale          Subcommand = "scale"
	Set            Subcommand = "set"
	Taint          Subcommand = "taint"
	Top            Subcommand = "top"
	Uncordon       Subcommand = "uncordon"
	Unknown        Subcommand = ""
	Version        Subcommand = "version"
	Wait           Subcommand = "wait"

	// oc (OpenShift CLI) specific subcommands
	Rsh Subcommand = "rsh"
)

func InspectSubcommand(cmdArgs []string, pluginHandler PluginHandler) (Subcommand, bool) {
	if len(cmdArgs) == 0 {
		return Unknown, false
	}
	cmd := cmdArgs[0]
	switch Subcommand(cmd) {
	case
		APIResources,
		APIVersions,
		Annotate,
		Apply,
		Attach,
		Auth,
		Autoscale,
		Certificate,
		ClusterInfo,
		Complete,
		CompleteNoDesc,
		Completion,
		Config,
		Convert,
		Cordon,
		Cp,
		Create,
		Debug,
		Delete,
		Describe,
		Diff,
		Drain,
		Edit,
		Events,
		Exec,
		Explain,
		Expose,
		Get,
		Kustomize,
		Label,
		Logs,
		Options,
		Patch,
		Plugin,
		PortForward,
		Proxy,
		Replace,
		Rollout,
		Rsh,
		Run,
		Scale,
		Set,
		Taint,
		Top,
		Uncordon,
		Version,
		Wait:
		return Subcommand(cmd), true
	default:
		if IsPlugin(cmdArgs, pluginHandler) {
			return KubectlPlugin, true
		}
		return Unknown, false
	}
}

func InspectSubcommandArgs(args []string, info *SubcommandInfo) {
	i := 0
	for i < len(args) {
		if args[i] == "--" {
			info.SubcommandArgs = append(info.SubcommandArgs, args[i+1:]...)
			break
		}

		arg, _, skip := parseArgFlag(args[i:])
		i += skip
		if !strings.HasPrefix(arg, "-") && skip == 1 {
			info.SubcommandArgs = append(info.SubcommandArgs, arg)
		}
	}
}

func CollectCommandlineOptions(args []string, info *SubcommandInfo) {
	i := 0
	for i < len(args) {
		// Stop parsing flags after "--", such as in "kubectl exec my-pod -- bash"
		if args[i] == "--" {
			break
		}
		flag, value, skip := parseArgFlag(args[i:])
		i += skip
		switch flag {
		case "--output", "-o":
			info.Output = ParseOutput(value)
		case "--client":
			info.Client = value != "false"
		case "--no-headers":
			info.NoHeader = true
		case "-w", "--watch", "--watch-only":
			info.Watch = true
		case "-f", "--follow":
			info.Follow = true
		case "--recursive":
			info.Recursive = value != "false"
		case "-h", "--help":
			info.Help = value != "false"
		case "--list":
			info.List = value != "false"
		}
	}
}

func parseArgFlag(args []string) (flag, value string, skip int) {
	if len(args) == 0 {
		return "", "", 0
	}
	arg := args[0]
	if strings.HasPrefix(arg, "--") {
		if flag, value, ok := strings.Cut(arg, "="); ok {
			// --output=wide
			return flag, value, 1
		}
		if len(args) > 1 {
			// --output wide
			return arg, args[1], 2
		}
	} else if strings.HasPrefix(arg, "-") && len(arg) >= 2 {
		if flag, value, ok := strings.Cut(arg, "="); ok {
			// -o=wide
			return flag, value, 1
		}
		if len(arg) > 2 {
			// -owide
			return arg[:2], arg[2:], 1
		}
		if len(arg) == 2 && len(args) > 1 {
			// -o wide
			return arg, args[1], 2
		}
	}
	return arg, "", 1
}

func InspectSubcommandInfo(args []string, pluginHandler PluginHandler) *SubcommandInfo {
	ret := &SubcommandInfo{}

	CollectCommandlineOptions(args, ret)

	for i, arg := range args {
		// Stop parsing args after "--", such as in "kubectl exec my-pod -- bash"
		if arg == "--" {
			break
		}

		cmd, ok := InspectSubcommand(args[i:], pluginHandler)
		if !ok {
			continue
		}

		ret.Subcommand = cmd
		InspectSubcommandArgs(args[i+1:], ret)
		return ret
	}

	// if subcommand is not found (e.g. kubecolor --help or just "kubecolor"),
	// it is treated as help because kubectl shows help for such input
	ret.Help = true

	return ret
}

func (sci *SubcommandInfo) SupportsPager() bool {
	if sci.Help {
		return false
	}
	switch sci.Subcommand {
	case Get:
		return !sci.Watch
	case Logs:
		return !sci.Follow
	case Describe,
		Explain,
		APIResources,
		APIVersions,
		Config:
		return true
	}
	return false
}

func (sci *SubcommandInfo) SupportsColoring() bool {
	switch sci.Subcommand {
	case Attach,
		Debug,
		Edit,
		Exec,
		Plugin,
		Proxy,
		Run,
		Wait:
		return sci.Help

	case KubectlPlugin,
		Complete, CompleteNoDesc:
		return false

	// oc (OpenShift CLI) specific subcommands
	case Rsh:
		return sci.Help

	// By default, all of our commands supports coloring
	default:
		return true
	}
}
