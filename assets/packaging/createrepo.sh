#!/usr/bin/env bash

# "createrepo" is the Fedora packaging tool used to create rpm/dnf/yum repositories

set -euo pipefail

createrepo_exec() {
    if command -v createrepo >/dev/null; then
        echo "# Running createrepo on your computer" >&2
        echo "# \$ createrepo" "$@" >&2
        createrepo "$@"
        return
    fi
    local docker=""
    if command -v podman >/dev/null; then
        docker=podman
    elif command -v docker >/dev/null; then
        docker=docker
    else
        echo "must have createrepo, docker, or podman installed"
        exit 1
    fi

    if [[ -z "$($docker images -q createrepo 2>/dev/null)" ]]; then
        echo "# Building image with createrepo" >&2
        $docker build . -t createrepo -f - <<EOF
FROM fedora:41
RUN dnf install -y createrepo_c
ENTRYPOINT ["createrepo"]
EOF
    fi

    echo "# Running createrepo inside Docker" >&2
    echo "# \$ $docker run --rm -it createrepo" "$@" >&2
    $docker run --rm -it -v "$PWD":/opt/src -w /opt/src createrepo "$@"
    echo >&2
}

createrepo_exec "$@"
