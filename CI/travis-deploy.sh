#!/bin/bash

set -ex

docker login --username "$DOCKERUSER" --password "$DOCKERPWD"

# Build images
docker build -t dangawne/points-are-bad/api-client -f api-client/Dockerfile .
docker built -t dangawne/points-are-bad/email-server -f backend/email-server/Dockerfile .
docker build -t dangawne/points-are-bad/fleetcommand-agent -f deploy/Dockerfile_fleetcommand_agent .

# Push images
docker push dangawne/points-are-bad/backend
docker push dangawne/points-are-bad/email-server
docker push dangawne/points-are-bad/fleetcommand-agent