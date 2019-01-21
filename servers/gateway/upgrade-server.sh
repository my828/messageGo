#!/bin/bash
docker pull my828/gateway
docker rm -f gateway 
export TLSCERT=/etc/letsencrypt/live/api.turtlemaster.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.turtlemaster.me/privkey.pem
docker run -d --name gateway -p 443:443 -e TLSCERT=$TLSCERT -e TLSKEY=$TLSKEY  -v /etc/letsencryt:/etc/letsencrypt:ro my828/gateway
