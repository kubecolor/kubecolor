#!/usr/bin/env bash

# Script that takes a list of files and prints them with different display names
# using the "${path}#${name}" format accepted by the "gh release create"
# & "gh release upload" commands.
#
# From:
#   dist/kubecolor_0.5.0_linux_amd64.deb
#   dist/kubecolor_0.5.0_linux_amd64.tar.gz
# To:
#   dist/kubecolor_0.5.0_linux_amd64.deb#kubecolor 0.5.0 linux amd64 deb
#   dist/kubecolor_0.5.0_linux_amd64.tar.gz#kubecolor 0.5.0 linux amd64

set -euo pipefail

if [ $# -eq 0 ]; then
  echo "usage: assets/packaging/label-assets.sh dist/kubecolor_*" >&2
  exit 1
fi

for asset in "$@"; do
  label="$(basename "$asset")"
  label="${label%.*}"
  label="${label%.tar}"
  label="$(echo "$label" | tr '_' ' ')"
  case "$asset" in
  *.msi ) label="${label} installer" ;;
  *.deb ) label="${label} deb" ;;
  *.rpm ) label="${label} RPM" ;;
  esac
  printf '"%s#%s"\n' "$asset" "$label"
done
