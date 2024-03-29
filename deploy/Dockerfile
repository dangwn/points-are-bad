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
RUN python3 -m venv /home/builduser/venv
ENV PATH="/home/builduser/venv/bin:$PATH"

# Clear cache
RUN apt-get clean && rm -rf /var/lib/apt/lists/*

# Install requirements
COPY ./fleetcommand/requirements.txt requirements.txt
RUN pip3 install --no-cache-dir wheel
RUN pip3 install --no-cache-dir --upgrade pip
RUN pip3 install --no-cache-dir -r requirements.txt

# Run image
FROM ubuntu:20.04 AS run-image

ARG DEBIAN_FRONTEND=noninterative

RUN apt-get -y update
RUN apt-get install --no-install-recommends -y python3.7 python3-venv curl libcurl4
 
# Install helm
RUN curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
RUN chmod 700 get_helm.sh
RUN ./get_helm.sh
RUN yes | rm get_helm.sh
RUN apt-get remove -y curl

RUN apt-get clean && rm -rf /var/lib/apt/lists/*

# Create build user
RUN useradd --create-home builduser
COPY --from=build-image /home/builduser/venv /home/builduser/venv

RUN mkdir /home/builduser/fleetcommand
WORKDIR /home/builduser/fleetcommand
COPY --chown=builduser:builduser ./fleetcommand .

# Change User
USER builduser
ENV PATH="/home/builduser/venv/bin:$PATH"

ENTRYPOINT [ "python", "-m", "agent" ]