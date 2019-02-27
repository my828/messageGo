#!/bin/bash
GOOS=linux go build 
docker build -t my828/messaging .
go clean
