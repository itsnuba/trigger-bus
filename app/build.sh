#!/bin/bash
set -e

# kadang harus sudo, https://github.com/docker/buildx/issues/415#issuecomment-761857828

DOCKERTAG=latest
if [ "$1" = "dev" ]; then
  DOCKERTAG=dev
fi

if [ "$2" != "skip-builder" ]; then
  docker build -f dockerfile.builder -t nabihubadah/trigger_bus-builder \
    .
fi

docker build -t nabihubadah/trigger_bus:$DOCKERTAG .
