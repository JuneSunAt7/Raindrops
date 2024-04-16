#!/bin/sh

GOOS=windows GOARCH=amd64 go build -o certs key_gens/generate_ssl.go

echo "Press enter to continue";
read name;
