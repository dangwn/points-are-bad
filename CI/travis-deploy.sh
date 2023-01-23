#!/bin/bash

set -ex

docker login --username "$DOCKERUSER" --password "$DOCKERPWD"

# Build images
docker build -t dangawne/points-are-bad-backend -f docker/Dockerfile_backend .
docker build -t dangawne/points-are-bad-fleetcommand-agent -f docker/Dockerfile_fleetcommand_agent .
# docker build -t dangawne/points-are-bad-web-ui -f docker/Dockerfile_web_ui .

# Push images
docker push dangawne/points-are-bad-backend
docker push dangawne/points-are-bad-fleetcommand-agent
# docker push dangawne/points-are-bad-web-ui