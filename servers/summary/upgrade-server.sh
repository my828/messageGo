#!/bin/bash
docker rm -f summary
docker pull my828/summary

# for messaging microservice
docker run -d \
--network auth \
-e ADDR=summary:80 \
--name summary \
my828/summary