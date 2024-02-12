---
title: Reference
description: HOw to use Kubecolor
---

## Installation

### Homebrew

![GitHub Release](https://img.shields.io/github/v/release/kubecolor/kubecolor?display_name=tag&label=Homebrew&color=4cc61f)

```sh
brew install kubecolor/tap/kubecolor
```

### Scoop

![Scoop Version](https://img.shields.io/scoop/v/kubecolor?label=Scoop%2FMain&color=4cc61f)

```sh
scoop install kubecolor
```

### Nix

[![nixpkgs unstable package](https://repology.org/badge/version-for-repo/nix_unstable/kubecolor.svg)](https://repology.org/project/kubecolor/versions)

```sh
nix-shell -p kubecolor
```

### AUR (Arch User Repositories)

[![AUR package](https://repology.org/badge/version-for-repo/aur/kubecolor.svg)](https://repology.org/project/kubecolor/versions)

```sh
yay -Syu kubecolor
```

### Download binary via GitHub release

Go to [Release page](https://github.com/kubecolor/kubecolor/releases) then download the binary which fits your environment.

### Compile from source

Requires Go 1.21 (or later)

```sh
go install github.com/kubecolor/kubecolor@latest
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
