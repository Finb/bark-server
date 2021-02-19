#!/usr/bin/env bash

set -e

if [ "$(docker version -f '{{.Server.Experimental}}')" != "true" ]; then
  echo "Docker daemon is not enabled experimental features."
  echo "Please follow the documentation to enable experimental features:"
  echo "  https://github.com/docker/cli/blob/master/experimental/README.md#use-docker-experimental"
  exit 1
fi

if [ ! -f ~/.docker/cli-plugins/docker-buildx ]; then
  echo "The buildx command is not installed."
  echo "Please follow the documentation to install buildx command:"
  echo "  https://github.com/docker/buildx/blob/master/README.md#installing"
  exit 1
fi

docker run --privileged --rm tonistiigi/binfmt --install all
docker buildx create --name bark-server --driver docker-container
docker buildx use bark-server
docker buildx build --platform linux/arm,linux/arm64,linux/386,linux/amd64 -t finab/bark-server:${BUILD_VERSION} -f deploy/Dockerfile --push .
docker buildx build --platform linux/arm,linux/arm64,linux/386,linux/amd64 -t finab/bark-server -f deploy/Dockerfile --push .
