# Points Are Bad

This is the repository containing the source code for [Points are Bad](https://pointsarebad.com).
---

This website runs a NextJS (in typescript) framework, with FastAPI backend, hosted inside an EKS cluster, with redis and postgres.

Contents:
- backend: A FastAPI-based REST API, which interacts with a postgresql database using sqlalchemy and alembic for the database migrations. User authentication happens using (unique and verified) emails and passwords, with access, refresh and csrf tokens (fastapi_csrf_protect) to authenticate users, with python-jose for password cryptography and validation.
- frontend: A NextJS frontend, written in typescript, which interacts with the FastAPI backend. The frontend uses react query for fetching data from the backend. The frontend will request that the user will accept the necessary cookies required for user authentication