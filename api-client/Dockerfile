FROM ubuntu:20.04 AS build-image

ARG DEBIAN_FRONTEND=noninterative

RUN apt-get -y update && apt-get upgrade -y
RUN apt-get install -y curl build-essential

# Install python
RUN apt-get install --no-install-recommends -y python3.7 \
    python3-dev \
    python3-venv \
    python3-pip \
    python3-wheel 

# Create virtual env
RUN python3 -m venv /home/apiuser/venv
ENV PATH="/home/apiuser/venv/bin:$PATH"

# Clear cache
RUN apt-get clean && rm -rf /var/lib/apt/lists/*

# Install requirements
COPY ./api-client/requirements.txt requirements.txt
RUN pip3 install --no-cache-dir wheel
RUN pip3 install --no-cache-dir --upgrade pip
RUN pip3 install --no-cache-dir -r requirements.txt

# Install rabbitmq-client
COPY ./backend-clients/tooling/rabbitmq-client ${HOME}/rabbitmq-client
WORKDIR ${HOME}/rabbitmq-client
RUN make build && make install && make clean 

# Run image
FROM ubuntu:20.04 AS run-image

ARG DEBIAN_FRONTEND=noninterative

RUN apt-get -y update
RUN apt-get install --no-install-recommends -y python3.7 python3-venv

RUN apt-get clean && rm -rf /var/lib/apt/lists/*

# Create API user
RUN useradd --create-home apiuser
COPY --from=build-image /home/apiuser/venv /home/apiuser/venv

RUN mkdir /home/apiuser/api
WORKDIR /home/apiuser/api
COPY --chown=apiuser:apiuser ./api-client/src .


# Change User
USER apiuser
ENV PATH="/home/apiuser/venv/bin:$PATH"

# Run API 
EXPOSE 8000