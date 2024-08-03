# Contributing

## Getting help

If you need any help, then we're happy to aid you in the GitHub issues.
Or if you prefer a more chat-like service, then we're also available
over at the Cloud-Native Slack workspace in the [#kubecolor](https://cloud-native.slack.com/archives/C06A9JNNHAR)
channel.

You can join the Slack workspace here: <https://communityinviter.com/apps/cloud-native/cncf>

## Prerequisites

- [Go](https://go.dev/) 1.21 (or later)
- GNUmake (the `make` command, used to run some steps in our [`Makefile`](./Makefile))

## Commit and branch policy

We don't really have a policy. Call the commits whatever you want.
Use `fix: ...` prefix if you want, or just name the commits what comes to mind.

Hell, maybe even use <https://whatthecommit.com>:

```bash
alias whatthecommit='git commit -m "$(curl -sSf https://whatthecommit.com/index.txt)"'
```

Preferably not nonsence messages though.

## Generating files

Note that some files are generated in this repo. To run them, use `make`:

```bash
# Regenerate ./config-schema.json
make config-schema.json

# Run integration test corpus, found in ./test/corpus/*.txt
make corpus

# Regenerate test results in integration test corpus
make corpus-update

# Regenerate screenshots used in README.md
make docs

# Generate configs (you currently have to copy-paste the results)
go run ./internal/cmd/configdoc
```
