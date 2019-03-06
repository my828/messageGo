#!/bin/bash
export TLSCERT=/etc/letsencrypt/live/api.turtlemaster.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.turtlemaster.me/privkey.pem
export MYSQL_ROOT_PASSWORD=$(openssl rand -base64 18)
export SESSIONKEY="sessionkey"
export REDISADDR="session:6379"
export DSN="root:$MYSQL_ROOT_PASSWORD@tcp(users:3306)/userinfo"
export MESSAGEADDR="http://message:80"
export SUMMARYADDR="http://summary:80"
export QNAME="rabbit"

docker rm -f gateway 
docker rm -f users
docker rm -f session
docker rm -f message
docker rm -f summary
docker rm -f rabbit
docker rm -f mongodb

docker network rm auth 

# create network 
docker network create auth

docker image prune -f
docker container prune -f
docker volume prune -f

docker pull my828/gateway
docker pull my828/database
docker pull rabbitmq:3

# for rabbit
docker run -d \
--network auth \
--name rabbit \
rabbitmq:3

# for redis
docker run -d \
--name session \
--network auth \
redis

# for mysql
docker run -d \
--network auth \
--name users \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_DATABASE=userinfo \
my828/database

sleep 20

docker run -d \
--name gateway \
--network auth \
-p 443:443 \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
-e DSN=$DSN \
-e REDISADDR=$REDISADDR \
-e SUMMARYADDR=$SUMMARYADDR \
-e MESSAGEADDR=$MESSAGEADDR \
-e QNAME=$QNAME \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
my828/gateway
