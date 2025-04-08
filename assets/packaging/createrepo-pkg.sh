#!/usr/bin/env bash

# Script that packages the output from goreleaser into an rpm repo.
#
# Everything only runs locally on local directories and only moves
# some files around in the dist directory. No network request.
# So it's safe to try out this command locally.

set -euo pipefail

dir="$(dirname "$0")"

# Clear out the destination
rm -rfv site/packages/rpm
mkdir -pv site/packages/rpm

# Create the repo metadata using "createrepo"
"$dir/createrepo.sh" --verbose dist --outputdir site/packages/rpm

# Copy over repo initialization file and rpm binaries
cp -v "$dir/kubecolor.repo" site/packages/rpm
cp -v dist/*.rpm site/packages/rpm

# TODO: Move to separate file
gpg --armor --detach-sign site/packages/rpm/repodata/repomd.xml
