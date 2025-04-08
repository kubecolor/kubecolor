#!/usr/bin/env bash

# Script that packages the output from goreleaser into a deb repo.
#
# Everything only runs locally on local directories and only moves
# some files around in the dist directory. No network request.
# So it's safe to try out this command locally.

set -euo pipefail

dir="$(dirname "$0")"

for file in dist/*.deb; do
    # The --confdir needs to contain the "distributions" file
    # "+b" means "$PWD"
    "$dir/reprepro.sh" \
        --confdir="+b/assets/packaging" \
        --outdir="+b/dist/deb" \
        --dbdir="+b/dist/deb/db" \
        includedeb stable "$file"
done
