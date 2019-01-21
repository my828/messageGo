#!/bin/bash
docker pull my828/summary2
export TLSCERT=/etc/letsencrypt/live/api.turtlemaster.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.turtlemaster.me/privkey.pem
docker run -d --name summary -p 80:80 -p 443:443 -v /etc/letsencryt:/etc/letsencrypt:ro my828/summary2
