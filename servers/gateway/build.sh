#!/bin/bash
GOOS=linux go build 
docker build -t my828/gateway .
docker build -t my828/database ../db
go clean

