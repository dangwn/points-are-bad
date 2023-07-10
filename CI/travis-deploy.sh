#!/bin/bash

set -ex

docker login --username "$DOCKERUSER" --password "$DOCKERPWD"

# Build images
docker build -t dangawne/points-are-bad/api-client -f backend/api-client/Dockerfile backend/api-client
docker build -t dangawne/points-are-bad/db-migrations -f backend/db-migrations/Dockerfile .
docker build -t dangawne/points-are-bad/email-server -f backend/email-server/Dockerfile backend/email-server
docker build -t dangawne/points-are-bad/fleetcommand-agent -f deploy/Dockerfile .

# Push images
docker push dangawne/points-are-bad/api-client
docker push dangawne/points-are-bad/db-migrations
docker push dangawne/points-are-bad/email-server
docker push dangawne/points-are-bad/fleetcommand-agent