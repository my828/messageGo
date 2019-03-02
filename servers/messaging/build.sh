#!/bin/bash
GOOS=linux go build 
docker build -t my828/message .
go clean
