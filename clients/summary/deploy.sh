#!/bin/bash
sh build.sh
docker push my828/summary2
ssh -oStrictHostKeyChecking=no ec2-user@ec2-35-162-127-249.us-west-2.compute.amazonaws.com 'bash -s' < upgrade-server.sh 