#!/bin/bash
export TLSCERT=/etc/letsencrypt/live/api.turtlemaster.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.turtlemaster.me/privkey.pem

docker rm -f gateway 
docker pull my828/gateway


docker run -d --name gateway -p 443:443 -e TLSCERT=$TLSCERT -e TLSKEY=$TLSKEY  -v /etc/letsencrypt:/etc/letsencrypt:ro my828/gateway
