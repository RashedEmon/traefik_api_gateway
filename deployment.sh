#!/bin/bash

image_name="traefik"
container_name="gateway"

# command to stop running container
container_stop_msg=$(docker stop "$container_name")
echo "container '$container_stop_msg' has stopped"

if docker ps -a --format '{{.Names}}' | grep -q "^$container_name$"; then
    docker rm -f $(docker ps -aqf "name=$container_name")
    echo "Containers with the name '$container_name' removed."
fi

# Check if the Docker image exists
if docker image inspect "$container_name" &> /dev/null; then
    docker image rm -f "$container_name"
    echo "Image '$container_name' deleted."
fi

# docker build
docker build -t "$image_name" .

# kill port if used
fuser --kill 4050/tcp
fuser --kill 3050/tcp

if ! [ -d /var/log/traefik/access_logs ]; then
    sudo mkdir -p /var/log/traefik/access_logs
fi

if ! [ -d /var/log/traefik/error_logs ]; then
    sudo mkdir -p /var/log/traefik/error_logs
fi

# run container
docker run -d -p 4050:8080 -p 3050:80 \
-v $PWD/traefik/traefik.yml:/etc/traefik/traefik.yml \
-v $PWD/traefik/dynamic_config.yml:/etc/traefik/dynamic_config.yml \
-v /var/log/traefik/access_logs:/var/log/traefik/access_logs \
-v /var/log/traefik/error_logs:/var/log/traefik/error_logs \
--restart=unless-stopped \
--name "$container_name" "$image_name"
