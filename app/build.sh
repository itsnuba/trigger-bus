#!/bin/bash
set -e

# harus sudo su, https://github.com/docker/buildx/issues/415#issuecomment-761857828

DOCKERTAG=latest
if [ "$1" = "dev" ]; then
  DOCKERTAG=dev
fi

if [ "$2" != "skip-builder" ]; then
  sudo docker build -f dockerfile.builder -t jala/v2/ms/lead-builder \
    --build-arg CONFIGURATION=$DOCKERBUILDCONF \
    --build-arg SSH_PRIVATE_KEY="$(< ~/.ssh/id_ed25519)" \
    .
fi

sudo docker build -t jala/v2/ms/lead:$DOCKERTAG .
