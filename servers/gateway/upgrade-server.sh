#!/bin/bash
docker pull my828/gateway
export TLSCERT=/etc/letsencrypt/live/api.turtlemaster.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.turtlemaster.me/privkey.pem
docker run -d --name gateway -p 443:443 -v /etc/letsencryt:/etc/letsencrypt:ro -e TLSCERT=$TLSCERT -e TLSKEY=$TLSKEY my828/gateway

