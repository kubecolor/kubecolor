#!/usr/bin/env bash

# Script that signs the output from goreleaser.
#
# The signatures are written to the rpm files,
# so this needs to be run before any "createrepo" tasks.

set -euo pipefail

dir="$(dirname "$0")"

"$dir/rpmsign.sh" --load=assets/packaging/rpmmacros --addsign dist/*.rpm
