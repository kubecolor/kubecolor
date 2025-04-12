#!/usr/bin/env bash

# Script that packages the output from goreleaser into a deb repo.
#
# Everything only runs locally on local directories and only moves
# some files around in the dist directory. No network request.
# So it's safe to try out this command locally.

set -euo pipefail

dir="$(dirname "$0")"

# Clear out the destination
rm -rfv packages/deb
mkdir -pv packages/deb

# Add the deb binaries to the repository one by one
for file in dist/*.deb; do
    # The --confdir needs to contain the "distributions" file
    # "+b" means "$PWD"
    "$dir/reprepro.sh" \
        --confdir="+b/assets/packaging" \
        --outdir="+b/packages/deb" \
        --dbdir="+b/packages/deb/db" \
        includedeb stable "$file"
done
