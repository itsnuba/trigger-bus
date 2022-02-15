#!/bin/bash
set -e

DOCKERTAG=latest
if [ "$1" = "dev" ]; then
  DOCKERTAG=dev
fi

docker run --rm --name trigger_bus -p 8900:8080 -it \
  -e "MONGO_URI=mongodb://host.docker.internal:27017" \
  -e "MONGO_DB=trigger-bus" \
  -e "DEBUG=0" \
  nabihubadah/trigger_bus:$DOCKERTAG
