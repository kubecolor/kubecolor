#!/usr/bin/env bash

# Script that just saves the latest version to a file so we can use it
# in installation scripts.

set -euo pipefail

mkdir -pv site/packages
echo -n 'version: '; jq '.version' dist/metadata.json -r | tee site/packages/version

# Debian uses "~" as a delimiter between the version and the "suffix"
# but still allows "-" in the suffix.
# Example: "0.5.0-foo-bar" -> "0.5.0~foo-bar"
mkdir -pv site/packages/deb
echo -n 'deb version: '; jq '.version | sub("(?<x>[^-]*)-(?<y>.*)"; "\(.x)~\(.y)")' dist/metadata.json -r | tee site/packages/deb/version
