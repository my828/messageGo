#!/bin/bash
export PORT="message:80"
export QNAME="rabbit"

docker rm -f message
docker rm -f mongodb

docker pull my828/message

docker run -d \
--network auth \
--name mongodb \
mongo


# for messaging microservice
docker run -d \
--network auth \
-e ADDR=$PORT \
-e QNAME=$QNAME \
--name message \
my828/message