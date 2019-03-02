#!/bin/bash
export PORT="message:80"
docker rm -f message
docker pull my828/message

docker run -d \
--network auth \
--name mongodb \
mongo

# for messaging microservice
docker run -d \
--network auth \
-e ADDR=$PORT \
--name message \
my828/message