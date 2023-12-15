<img alt="Kubecolor logo" src="./assets/img/Kubecolor_Logo_large.png" />

# kubecolor

KubeColor is a `kubectl` replacement used to add colors to your kubectl output.

![test](https://github.com/kubecolor/kubecolor/workflows/test/badge.svg?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/kubecolor/kubecolor)](https://goreportcard.com/report/github.com/kubecolor/kubecolor)
[![codecov](https://codecov.io/gh/kubecolor/kubecolor/branch/main/graph/badge.svg?token=k6ysAa5ghD)](https://codecov.io/gh/kubecolor/kubecolor/)


* get pods

![image](https://user-images.githubusercontent.com/60682957/95733375-04929680-0cbd-11eb-82f3-adbcfecf4a3e.png)

* describe pods

![image](https://user-images.githubusercontent.com/60682957/95733389-08beb400-0cbd-11eb-983b-cf5138277fe3.png)

* something wrong

![image](https://user-images.githubusercontent.com/60682957/95733397-0a887780-0cbd-11eb-8875-bb1000e0e597.png)

* You can change color theme for light-backgrounded environment

![image](https://user-images.githubusercontent.com/60682957/95733403-0c523b00-0cbd-11eb-9ff9-abc5469e97ca.png)

## What's this?

kubecolor colorizes your `kubectl` command output and does nothing else.
kubecolor internally calls `kubectl` command and try to colorizes the output so
you can use kubecolor as a complete alternative of kubectl. It means you can write this in your .bash_profile:

```sh
alias kubectl="kubecolor"
```
If you use your .bash_profile on more than one computer (e.g. synced via git) that might not all have `kubecolor`
installed, you can avoid breaking `kubectl` like so:

```sh
command -v kubecolor >/dev/null 2>&1 && alias kubectl="kubecolor"
```

For now, not all subcommands are supported and some of them might be in the future. If you want something to be supported by kubecolor, you can open an issue in this GitHub repo.
Even if what you want to do is not supported by kubecolor right now, kubecolor still can just show `kubectl` output without any decorations,
so you don't need to switch kubecolor and kubectl but you always can use kubecolor.

Additionally, if `kubectl` resulted an error, kubecolor just shows the error message in red or yellow.

**Please keep in mind that kubecolor has runtime overhead compared to kubectl because it runs kubectl internally and parses the output, then colorizes it.**

## Installation

### Download binary via GitHub release

Go to [Release page](https://github.com/kubecolor/kubecolor/releases) then download the binary which fits your environment.

### Mac and Linux users via Homebrew

```sh
brew install kubecolor/tap/kubecolor
```

### Manually via go command

*Note: if you install kubecolor via go command without setting the `ldflags`, --kubecolor-version  might return an unset value*

```sh
VERSION=latest
go install -ldflags="-X main.Version=${VERSION}" github.com/kubecolor/kubecolor/cmd/kubecolor@${VERSION}
```

You can install a specific version like:

```sh
VERSION=v0.0.21
go install -ldflags="-X main.Version=${VERSION}" github.com/kubecolor/kubecolor/cmd/kubecolor@${VERSION}
```

If you are not using module mode (old Go version < 1.11 or command above doesn't work), try below:

```sh
go get -u github.com/kubecolor/kubecolor/cmd/kubecolor
```

## Usage

kubecolor understands every subcommands and options which are available for `kubectl`. What you have to do is just using `kubecolor`
instead of `kubectl` like:

```sh
kubecolor --context=your_context get pods -o json
```

If you want to make the colorized kubectl default on your shell, just add this line into your shell configuration file:

```sh
alias kubectl="kubecolor"
```

### Dynamic tty support

When the kubecolor output tty is not standard output, it automatically disables the colorization.
For example, if you are running `kubecolor get pods > result.txt` or `kubecolor get pods | grep xxx`, the output will be passed through to file or another command, so colorization is not applied.
You can force kubecolor do colorization at such cases by passing `--force-colors` flag. See the upcoming section for more details.

### Flags

Available flags for kubecolor. When you pass them, kubecolor will understand them but these flags won't be passed to kubectl.

* `--kubecolor-version`

Prints the version of kubecolor (not kubectl one).

* `--light-background`

When your terminal's background color is something light (e.g white), default color preset might look too bright and not readable.
If so, specify `--light-background` as a command line argument. kubecolor will use a color preset for light-backgrounded environment.

* `--force-colors`

By default, kubecolor never output the result in colors when the tty is not a terminal standard output.
If you want to force kubecolor to show the result in colors for non-terminal tty, you can specify this flag.
For example, when you want to pass kubecolor result to grep (`kubecolor get pods | grep pod_name`), this option is useful.

* `--plain`

When you don't want to colorize output, you can specify `--plain`. Kubecolor understands this option and outputs the result without colorizing.

### ENV Variables

* `KUBECOLOR_OBJ_FRESH`

When setting the variable `KUBECOLOR_OBJ_FRESH` to a duration, you can change the color of the object depending on its creation time.

Please see [Specify object fresh age threshold](#specify-object-fresh-age-threshold)

* `KUBECOLOR_FORCE_COLORS`

In addition to forcing colors with `--force-colors`, you can also do so by setting the environment variable `KUBECOLOR_FORCE_COLORS=true`.
You can use this environment variable to colorize output when you invoke kubecolor in the `watch` command (e.g. `watch kubecolor get pods`). 
Set the following alias:
```shell
alias watch='KUBECOLOR_FORCE_COLORS=true watch --color '
```

Be sure to include the space at the end to enable alias expansion (without this additional space, the command `watch kgp` would fail, for example). 

### Autocompletion

kubectl provides [autocompletion feature](https://kubernetes.io/docs/tasks/tools/install-kubectl/#enable-kubectl-autocompletion). If you are
already using it, you might have to configure it for kubecolor.
Please also refer to [kubectl official doc for kubectl autocomplete](https://kubernetes.io/docs/reference/kubectl/cheatsheet/#kubectl-autocomplete).

#### Bash

For Bash, configuring autocompletion requires adding following line in your shell config file.

```shell
# autocomplete for kubecolor
complete -o default -F __start_kubectl kubecolor
```

If you are using an alias like `k="kubecolor"`, then just change above like:

```shell
complete -o default -F __start_kubectl k
```

#### Zsh

For zsh make sure these lines are present in your zsh config file:
```shell
# get zsh complete kubectl
source <(kubectl completion zsh)
alias kubectl=kubecolor
# make completion work with kubecolor
compdef kubecolor=kubectl
```

#### fish

Fish completion is officially unsupported by `kubectl`, so it is unsupported by `kubecolor` as well.

However, there are 2 ways we can make them work. Please keep in mind these are a kind of "hack" and not officially supported.

1. Use [evanlucas/fish-kubectl-completions](https://github.com/evanlucas/fish-kubectl-completions) with `kubecolor`:

   * Install `kubectl` completions (https://github.com/evanlucas/fish-kubectl-completions):
      ```
      fisher install evanlucas/fish-kubectl-completions
      ```
   * Add the following function to your `config.fish` file:
      ```
      function kubectl
        kubecolor $argv
      end
      ```
2. Use [awinecki/fish-kubecolor-completions](https://github.com/awinecki/fish-kubecolor-completions)

   The first way will override `kubectl` command. If you wish to preserve both `kubectl` and `kubecolor` with completions, you need to copy [evanlucas/fish-kubectl-completions](https://github.com/evanlucas/fish-kubectl-completions) for the `kubecolor` command.

   For this purpose, you can use [awinecki/fish-kubecolor-completions](https://github.com/awinecki/fish-kubecolor-completions).


### Specify what command to execute as kubectl

Sometimes, you may want to specify which command to use as `kubectl` internally in kubecolor. For example, when you want to use a versioned-kubectl `kubectl.1.19`, you can do that by an environment variable:

```shell
KUBECTL_COMMAND="kubectl.1.19" kubecolor get po
```

When you don't set `KUBECTL_COMMAND`, then `kubectl` is used by default.

### Specify object fresh age threshold

When you want to figure out which object are recently created the fresh ones are colored green. It only works with `get` subcommand.
You can customize threshold level by using environment variable `KUBECOLOR_OBJ_FRESH` with value format [time.Duration](https://golang.org/pkg/time/#ParseDuration).

```shell
KUBECOLOR_OBJ_FRESH="1m" kubecolor get po
```

Default value is `0s`, it means is disabled.
## Supported kubectl version

Because kubecolor internally calls `kubectl` command, if you are using unsupported kubectl version, it's also not supported by kubecolor.
Kubernetes version support policy can be found in [official doc](https://kubernetes.io/docs/setup/release/version-skew-policy/).

## kubectl plugins

Even if kubectl supports [plugin](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/) feature, it is not supported by kubecolor, including [Krew](https://krew.sigs.k8s.io/).
It means you can use plugins from kubecolor (e.g. you can do `kubecolor plugin_name xxx`), but the result won't be colorized.

## Contributions

Always welcome. Just opening an issue should be also greatful.

## LICENSE

MIT

## Author

This project is a heavily modified version of the original KubeColor maintained at [https://github.com/hidetatz/kubecolor](https://github.com/hidetatz/kubecolor)

[@kubecolor](https://github.com/kubecolor)
