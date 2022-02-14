#!/bin/bash
set -e

DOCKERTAG=latest
if [ "$1" = "dev" ]; then
  DOCKERTAG=dev
fi

aws ecr get-login-password --region ap-southeast-1 | docker login --username AWS --password-stdin 352635696299.dkr.ecr.ap-southeast-1.amazonaws.com

docker tag jala/v2/ms/lead:$DOCKERTAG 352635696299.dkr.ecr.ap-southeast-1.amazonaws.com/jala/v2/ms/lead:$DOCKERTAG

docker push 352635696299.dkr.ecr.ap-southeast-1.amazonaws.com/jala/v2/ms/lead:$DOCKERTAG
