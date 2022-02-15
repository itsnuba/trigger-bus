#!/bin/bash
set -e

DOCKERTAG=latest
if [ "$1" = "dev" ]; then
  DOCKERTAG=dev
fi

docker tag nabihubadah/trigger_bus:$DOCKERTAG nabihubadah/trigger_bus:$DOCKERTAG

docker push nabihubadah/trigger_bus:$DOCKERTAG
