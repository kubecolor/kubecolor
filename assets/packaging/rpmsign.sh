#!/usr/bin/env bash

# "rpmsign" is the Fedora packaging tool used to sign rpm packages using GPG.

set -euo pipefail

rpmsign_exec() {
    if command -v rpmsign >/dev/null; then
        echo "# Running rpmsign on your computer" >&2
        echo "# \$ rpmsign" "$@" >&2
        rpmsign "$@"
        return
    fi
    local docker=""
    if command -v podman >/dev/null; then
        docker=podman
    elif command -v docker >/dev/null; then
        docker=docker
    else
        echo "must have rpmsign, docker, or podman installed"
        exit 1
    fi

    if [[ -z "$($docker images -q rpmsign 2>/dev/null)" ]]; then
        echo "# Building image with rpmsign" >&2
        $docker build . -t rpmsign -f - <<EOF
FROM fedora:41
RUN dnf install -y rpm-sign
ENTRYPOINT ["rpmsign"]
EOF
    fi

    echo "# Running rpmsign inside Docker" >&2
    echo "# \$ $docker run --rm -it rpmsign" "$@" >&2
    $docker run --rm -it -v "$PWD":/opt/src -v "${GNUPGHOME:-"$HOME"/.gnupg}":/root/.gnupg -w /opt/src rpmsign "$@"
    echo >&2
}

rpmsign_exec "$@"
