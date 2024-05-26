package kubectl

import (
	"strings"
)

type SubcommandInfo struct {
	Subcommand   Subcommand
	FormatOption FormatOption
	NoHeader     bool
	Watch        bool
	Follow       bool
	Help         bool
	Recursive    bool
	Client       bool
}

type FormatOption int

const (
	None FormatOption = iota
	Wide
	Json
	Yaml
)

type Subcommand string

const (
	Unknown          Subcommand = ""
	KubectlPlugin    Subcommand = "(plugin)"
	InternalComplete Subcommand = "__complete"

	APIResources Subcommand = "api-resources"
	APIVersions  Subcommand = "api-versions"
	Annotate     Subcommand = "annotate"
	Apply        Subcommand = "apply"
	Attach       Subcommand = "attach"
	Auth         Subcommand = "auth"
	Autoscale    Subcommand = "autoscale"
	Certificate  Subcommand = "certificate"
	ClusterInfo  Subcommand = "cluster-info"
	Completion   Subcommand = "completion"
	Config       Subcommand = "config"
	Convert      Subcommand = "convert"
	Cordon       Subcommand = "cordon"
	Cp           Subcommand = "cp"
	Create       Subcommand = "create"
	Ctx          Subcommand = "ctx"
	Debug        Subcommand = "debug"
	Delete       Subcommand = "delete"
	Describe     Subcommand = "describe"
	Diff         Subcommand = "diff"
	Drain        Subcommand = "drain"
	Edit         Subcommand = "edit"
	Events       Subcommand = "events"
	Exec         Subcommand = "exec"
	Explain      Subcommand = "explain"
	Expose       Subcommand = "expose"
	Get          Subcommand = "get"
	Kustomize    Subcommand = "kustomize"
	Label        Subcommand = "label"
	Logs         Subcommand = "logs"
	Ns           Subcommand = "ns"
	Options      Subcommand = "options"
	Patch        Subcommand = "patch"
	Plugin       Subcommand = "plugin"
	PortForward  Subcommand = "port-forward"
	Proxy        Subcommand = "proxy"
	Replace      Subcommand = "replace"
	Rollout      Subcommand = "rollout"
	Run          Subcommand = "run"
	Scale        Subcommand = "scale"
	Set          Subcommand = "set"
	Taint        Subcommand = "taint"
	Top          Subcommand = "top"
	Uncordon     Subcommand = "uncordon"
	Version      Subcommand = "version"
	Wait         Subcommand = "wait"

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
		Completion,
		Config,
		Convert,
		Cordon,
		Cp,
		Create,
		Ctx,
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
		Ns,
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
		// Catch __complete, __completeNoDesc, etc
		if strings.HasPrefix(cmd, "__complete") {
			return InternalComplete, true
		}

		if IsPlugin(cmdArgs, pluginHandler) {
			return KubectlPlugin, true
		}
		return Unknown, false
	}
}

func CollectCommandlineOptions(args []string, info *SubcommandInfo) {
	for i := range args {
		// Stop parsing flags after "--", such as in "kubectl exec my-pod -- bash"
		if args[i] == "--" {
			break
		}
		if strings.HasPrefix(args[i], "--output") {
			switch args[i] {
			case "--output=json":
				info.FormatOption = Json
			case "--output=yaml":
				info.FormatOption = Yaml
			case "--output=wide":
				info.FormatOption = Wide
			default:
				if len(args)-1 > i {
					formatOption := args[i+1]
					switch formatOption {
					case "json":
						info.FormatOption = Json
					case "yaml":
						info.FormatOption = Yaml
					case "wide":
						info.FormatOption = Wide
					default:
						// custom-columns, go-template, etc are currently not supported
					}
				}
			}
		} else if strings.HasPrefix(args[i], "-o") {
			switch args[i] {
			// both '-ojson' and '-o=json' works
			case "-ojson", "-o=json":
				info.FormatOption = Json
			case "-oyaml", "-o=yaml":
				info.FormatOption = Yaml
			case "-owide", "-o=wide":
				info.FormatOption = Wide
			default:
				// otherwise, look for next arg because '-o json' also works
				if len(args)-1 > i {
					formatOption := args[i+1]
					switch formatOption {
					case "json":
						info.FormatOption = Json
					case "yaml":
						info.FormatOption = Yaml
					case "wide":
						info.FormatOption = Wide
					default:
						// custom-columns, go-template, etc are currently not supported
					}
				}

			}
		} else if strings.HasPrefix(args[i], "--client") {
			switch args[i] {
			case "--client=true":
				info.Client = true
			case "--client=false":
				info.Client = false
			default:
				info.Client = true
			}
		} else if args[i] == "--no-headers" {
			info.NoHeader = true
		} else if args[i] == "-w" || args[i] == "--watch" || args[i] == "--watch-only" {
			info.Watch = true
		} else if args[i] == "-f" || args[i] == "--follow" {
			info.Follow = true
		} else if args[i] == "--recursive=true" || args[i] == "--recursive" {
			info.Recursive = true
		} else if args[i] == "-h" || args[i] == "--help" {
			info.Help = true
		}
	}
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
		Explain:
		return true
	}
	return false
}
