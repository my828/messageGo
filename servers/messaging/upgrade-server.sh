#!/bin/bash
export PORT="80"
docker pull my828/messaging

docker rm -f session

# for messaging microservice
docker run -d \
--network auth \
--name message \
my828/messaging