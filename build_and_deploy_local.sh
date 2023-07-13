#!/bin/bash

BUILD_TAG="latest"
BUILD_CONFIG="docker-deploy-local.yml"
RUN_MIGRATIONS_SCRIPT="backend/db-migrations/run_migrations.py"

help() {
    echo "Build dockerfiles and deploy them using docker swarm"
    echo "Flags: "
    echo "  -f | Path to docker deploy configuration yaml file (default 'docker-deploy-local.yml')"
    echo "  -m | Path to run migrations python script (default 'db-migrations/run_migrations.py')"
    echo "  -t | Tag for the PAB image builds referenced in docker deploy configuration file (default 'latest')"
    echo "  -h | Help"
}

while getopts "f:m:t:h" flag; do
    case "${flag}" in
        f) BUILD_CONFIG="${OPTARG}" ;;
        m) RUN_MIGRATIONS_SCRIPT="${OPTARG}" ;;
        t) BUILD_TAG="${OPTARG}" ;;
        h) help
            exit 0 ;;
        *) help
            exit 1 ;;
    esac
done

echo "Building points are bad images with tag '${BUILD_TAG}'..."
echo ""
echo "Building api client image..."
docker build -f backend/api-client/Dockerfile -t dangawne/points-are-bad/api-client:$BUILD_TAG --build-arg API_PORT=8020 backend/api-client
echo ""
echo "Building email server image..."
docker build -f backend/email-server/Dockerfile -t dangawne/points-are-bad/email-server:$BUILD_TAG backend/email-server
echo ""
echo "Building OpenAPI docs image..."
docker build -f backend/openapi/Dockerfile -t dangawne/points-are-bad/swagger-ui:$BUILD_TAG backend/openapi
echo ""
echo "Building migrations image..."
docker build -f backend/db-migrations/Dockerfile -t dangawne/points-are-bad/db-migrations:$BUILD_TAG backend/db-migrations
echo ""
echo "Images built successfully!"

echo ""
echo "Pulling external images..."
docker pull redis:7.0.10
docker pull postgres:10.17
docker pull adminer:latest
docker pull rabbitmq:3-management
echo ""
echo "External images pulled successfully!"

echo ""
echo "Deploying in docker swarm..."
docker swarm init
docker network create -d overlay --attachable pab_public
docker stack deploy -c ${BUILD_CONFIG} --with-registry-auth pab

echo ""
echo "Running db migrations"
sleep 10
docker run --network=pab_public dangawne/points-are-bad/db-migrations:latest

echo ""
echo "Done!"