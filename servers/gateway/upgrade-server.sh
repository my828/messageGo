#!/bin/bash
export TLSCERT=/etc/letsencrypt/live/api.turtlemaster.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.turtlemaster.me/privkey.pem
export MYSQL_ROOT_PASSWORD=$(openssl rand -base64 18)
export SESSIONKEY="sessionkey"
export REDISADDR="session:6379"
export DSN="root:$MYSQL_ROOT_PASSWORD@tcp(users:3306)/userinfo"

docker rm -f gateway 
docker rm -f users
docker rm -f session
docker network rm auth 


# docker network disconnect -f auth gateway
# docker network disconnect -f auth users
# docker network disconnect -f auth session

docker pull my828/gateway
docker pull my828/database


docker network create auth

docker run -d \
--name session \
--network auth \
redis

docker run -d \
--network auth \
--name users \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_DATABASE=userinfo \
my828/database

sleep 15

docker run -d \
--name gateway \
--network auth \
-p 443:443 \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
-e DSN=$DSN \
-e REDISADDR=$REDISADDR \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
my828/gateway
