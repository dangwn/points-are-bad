#!/bin/bash

BUILD_TAG="latest"
BUILD_CONFIG="docker-deploy-local.yml"
RUN_MIGRATIONS_SCRIPT="db-migrations/run_migrations.py"

help() {
    echo "Build dockerfiles and deploy them using docker swarm"
    echo "Flags: "
    echo "  -f | Path to docker deploy configuration yaml file (default 'docker-deploy-local.yml')"
    echo "  -m | Path to run migrations python script (default 'db-migrations/run_migrations.py')"
    echo "  -t | Tag for the PAB image builds referenced in docker deploy configuration file (default 'latest')"
    echo "  -h | Help"
}

while getopts "t:f:h" flag; do
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
echo "Pulling external images..."
docker pull redis:7.0.10
docker pull postgres:10.17
docker pull adminer:latest
docker pull rabbitmq:3-management

echo ""
echo "Deploying in docker swarm..."
docker swarm init
docker stack deploy -c ${BUILD_CONFIG} pab


echo ""
echo "Running db migrations"
sleep 10
python ${RUN_MIGRATIONS_SCRIPT}

echo ""
echo "Done!"