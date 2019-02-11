#!/bin/bash

docker rm -f summary2
docker pull my828/summary2

docker run -d -p 80:80 -p 443:443 -v /etc/letsencrypt:/etc/letsencrypt:ro my828/summary2
