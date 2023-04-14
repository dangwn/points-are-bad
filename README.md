# Points Are Bad

This is the repository containing the source code for [Points are Bad](https://pointsarebad.com).
---

This website runs a NextJS (in typescript) framework, with FastAPI backend, hosted inside an EKS cluster, with redis and postgres.

## Contents:
- <b>backend</b>: A FastAPI-based REST API, which interacts with a postgresql database using sqlalchemy and alembic for the database migrations. User authentication happens using (unique and verified) emails and passwords, with access, refresh and csrf tokens (fastapi_csrf_protect) to authenticate users, with python-jose for password cryptography and validation.
- <b>CI</b>: Assets for CI of the project. Currently only holds shell script for travis to build the docker images.
- <b>Deploy</b>: Assets for deploying the backend onto the EKS cluster. Contains helm charts for K8s resources. A shell script deploys a pod onto the K8s cluster, which installs the required helm charts with the given configuration.
- <b>Docker</b>: Dockerfiles for the fleetcommand agent and the backend API.
- <b>frontend</b>: A NextJS frontend, written in typescript, which interacts with the FastAPI backend. The frontend uses react query for fetching data from the backend. The frontend will request that the user will accept the necessary cookies required for user authentication