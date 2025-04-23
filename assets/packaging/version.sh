#!/usr/bin/env bash

# Script that just saves the latest version to a file so we can use it
# in installation scripts.

set -euo pipefail

: "${GITHUB_OUTPUT:=}"

function set-output() {
    local name="$1"
    local value="$2"

    echo "$name=$value"
    if [[ -n "$GITHUB_OUTPUT" ]]; then
        echo "$name=$value" >> "$GITHUB_OUTPUT"
    fi
}

if [[ ! -f dist/metadata.json ]]; then
    echo "Missing dist/metadata.json file"
    exit 1
fi

echo "metadata.json=$(cat dist/metadata.json)"

mkdir -pv packages
VERSION="$(jq '.version' dist/metadata.json -r | tee packages/version)"
set-output version "$VERSION"

# Debian uses "~" as a delimiter between the version and the "suffix"
# but still allows "-" in the suffix.
# Example: "0.5.0-foo-bar" -> "0.5.0~foo-bar"
mkdir -pv packages/deb
DEB_VERSION="$(jq '.version | sub("(?<x>[^-]*)-(?<y>.*)"; "\(.x)~\(.y)")' dist/metadata.json -r | tee packages/deb/version)"
set-output deb-version "$DEB_VERSION"
