#!/usr/bin/env bash

# "reprepro" is the Debian packaging tool used to create deb/apt repositories
# and sign it using GPG.

set -euo pipefail

reprepro_exec() {
    if command -v reprepro >/dev/null; then
        echo "# Running reprepro on your computer" >&2
        echo "# \$ reprepro" "$@" >&2
        reprepro "$@"
        return
    fi
    local docker=""
    if command -v podman >/dev/null; then
        docker=podman
    elif command -v docker >/dev/null; then
        docker=docker
    else
        echo "must have reprepro, docker, or podman installed"
        exit 1
    fi

    if [[ -z "$($docker images -q reprepro 2>/dev/null)" ]]; then
        echo "# Building image with reprepro" >&2
        $docker build . -t reprepro -f - <<EOF
FROM ubuntu:24.04
RUN apt update && apt install reprepro -y
ENTRYPOINT ["reprepro"]
EOF
    fi

    echo "# Running reprepro inside Docker" >&2
    echo "# \$ $docker run --rm reprepro" "$@" >&2
    $docker run --rm -v "$PWD":/opt/src -v "${GNUPGHOME:-"$HOME"/.gnupg}":/root/.gnupg -w /opt/src reprepro "$@"
    echo >&2
}

reprepro_exec "$@"
