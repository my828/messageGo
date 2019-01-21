#!/bin/bash
ssh ec2-user@ec2-35-162-127-249.us-west-2.compute.amazonaws.com 
GOOS=linux go build .
docker build -t my828/summary.
docker push my828/summary
docker rm -f summary 
docker pull my828/summary
export TLSCERT=/etc/letsencrypt/live/api.turtlemaster.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.turtlemaster.me/privkey.pem
docker run -d --name summary -p 443:443 -v /etc/letsencryt:/etc/letsencrypt:ro -e TLSCERT=$TLSCERT -e TLSKEY=$TLSKEY my828/summary 
