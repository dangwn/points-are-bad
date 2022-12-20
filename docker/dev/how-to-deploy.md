## Steps for deploying development stack
From root directory:
1. `docker pull postgres adminer`
2. `docker build -t pab-api -f docker/dev/Dockerfile_api .`
3. `docker build -t pab-migrations -f docker/dev/Dockerfile_api .`
4. Change `DB_HOST` environment variables in `api` and `migrations` in `pab_stack.yml` file to be private IP
5. `docker swarm init`
6. `docker stack deploy -c docker/dev/pab_stack.yml pab`

Then to leave: <br>
`docker swarm leave --force && docker system prune`