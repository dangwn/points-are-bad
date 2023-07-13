# Points Are Bad

This is the repository containing the source code for [Points are Bad](https://pointsarebad.com).
---

## Contents
 - <b>backend:</b> Source code for containers running in the k8s cluster, the backend of Points are Bad
    - <b>api-client:</b> The REST API which handles the request from the frontend
    - <b>email-server:</b> Server which handles all outgoing emails to users, either account verification or password resets
    - <b>db-migrations:</b> Python scripts defining SQL models and for running database migrations
    - <b>open-api:</b> Open API spec for api client with dockerfile for local deploy
 - <b>CI:</b> Script for building and pushing docker images using [Travis CI](https://travis-ci.com)
 - <b>deploy:</b> Deployment code and helm charts for deploying k8s resources onto the cluster
 - <b>frontend:</b> Source code for the frontend web application

## Deploy locally
### Requirements:
 - [Docker](https://docs.docker.com/) >= 1.13 (must have [docker swarm](https://docs.docker.com/engine/swarm/) available)
 - An AWS SES user deployed, with at least one email address verified to send emails from ([SES docs](https://docs.aws.amazon.com/ses/latest/dg/send-email.html))

Run the `build_and_deploy_local.sh` script (add `-h` for help) from the root directory of the project.