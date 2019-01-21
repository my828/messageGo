#!/bin/bash
GOOS=linux go build 
docker build -t my828/summary2 .
go clean