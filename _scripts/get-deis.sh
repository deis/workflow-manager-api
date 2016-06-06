#!/usr/bin/env bash
#
# Download the deis client package and install to cwd
#

cd "$(dirname "$0")" || exit 1

curl -sSL http://deis.io/deis-cli/install-v2.sh | bash
