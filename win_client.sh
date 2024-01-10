#!/bin/sh

GOOS=linux GOARCH=amd64 go build -o admin cmd/data_mgmt/server_mgmt.go 

echo "Press enter to continue";
read name;
