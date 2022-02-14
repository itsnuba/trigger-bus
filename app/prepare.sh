#!/bin/bash
set -e

# pindah ke folder app
cd $(dirname "$0")

# install air
# https://github.com/cosmtrek/air
curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b ../bin
../bin/air -v
