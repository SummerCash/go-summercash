#!/bin/bash
# Note: this script was taken from the go-spacemesh repository.
# All credits for this script, validate_lint.sh, go to the spacemeshos development team.

BRANCH=$(git rev-parse --abbrev-ref HEAD) # Get branch

docker build -t go-summercash:$BRANCH . # Build docker image
echo $DOCKER_PASSWORD | docker login -u $DOCKER_USERNAME --password-stdin
docker tag go-summercash:$BRANCH summercash/go-summercash:$BRANCH
docker push summercash/go-summercash:$BRANCH