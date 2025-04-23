#!/usr/bin/env bash

# Script that signs the repomd.xml file using GPG.

set -euo pipefail

gpg --yes --armor --detach-sign packages/rpm/repodata/repomd.xml
