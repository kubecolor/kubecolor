#!/usr/bin/env bash

# Script that packages the output from goreleaser into an rpm repo.
#
# Everything only runs locally on local directories and only moves
# some files around in the dist directory. No network request.
# So it's safe to try out this command locally.

set -euo pipefail

dir="$(dirname "$0")"

"$dir/createrepo.sh" -v dist
