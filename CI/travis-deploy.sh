#!/bin/bash

set -ex

docker login --username "$DOCKERUSER" --password "$DOCKERPWD"

# Build images
docker build -t dangawne/points-are-bad-backend -f api-client/Dockerfile_backend .
docker build -t dangawne/points-are-bad-fleetcommand-agent -f deploy/Dockerfile_fleetcommand_agent .

# Push images
docker push dangawne/points-are-bad-backend
docker push dangawne/points-are-bad-fleetcommand-agent