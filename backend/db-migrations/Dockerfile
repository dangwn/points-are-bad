FROM ubuntu:20.04 AS build-image

ARG DEBIAN_FRONTEND=noninterative

RUN apt-get -y update && apt-get upgrade -y
RUN apt-get install -y curl

# Install python
RUN apt-get install --no-install-recommends -y python3.7 \
    python3-dev \
    python3-venv \
    python3-pip \
    python3-wheel 

# Create virtual env
RUN python3 -m venv /home/migrationsuser/venv
ENV PATH="/home/migrationsuser/venv/bin:$PATH"

# Clear cache
RUN apt-get clean && rm -rf /var/lib/apt/lists/*

# Install requirements
COPY backend/db-migrations/requirements.txt requirements.txt
RUN pip3 install --no-cache-dir wheel
RUN pip3 install --no-cache-dir --upgrade pip
RUN pip3 install --no-cache-dir -r requirements.txt

# Run image
FROM ubuntu:20.04 AS run-image

ARG DEBIAN_FRONTEND=noninterative

RUN apt-get -y update
RUN apt-get install --no-install-recommends -y python3.7 python3-venv

RUN apt-get clean && rm -rf /var/lib/apt/lists/*

# Create API user
RUN useradd --create-home migrationsuser
COPY --from=build-image /home/migrationsuser/venv /home/migrationsuser/venv

RUN mkdir /home/migrationsuser/api
WORKDIR /home/migrationsuser/api
COPY --chown=migrationsuser:migrationsuser backend/db-migrations/src .

# Change User
USER migrationsuser
ENV PATH="/home/migrationsuser/venv/bin:$PATH"