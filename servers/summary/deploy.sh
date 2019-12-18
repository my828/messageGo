#!/bin/bash
sh build.sh
docker push my828/summary
ssh -oStrictHostKeyChecking=no ec2-user@ec2-3-94-215-128.compute-1.amazonaws.com 'bash -s' < upgrade-server.sh 
