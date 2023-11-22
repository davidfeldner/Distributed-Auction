#!/bin/bash
go run ./server/server.go 8080
go run ./server/server.go 8081
go run ./server/server.go 8082
sleep 1
go run ./client/client.go servers.txt
go run ./client/client.go servers.txt