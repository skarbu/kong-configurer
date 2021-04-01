#!/bin/bash
export GOOS="linux"
export GOARCH="amd64"
go build -o ./linux-amd64/kong-configurer ../kong-configurer.go
export GOOS="windows"
export GOARCH="amd64"
go build -o ./windows-amd64/kong-configurer.exe ../kong-configurer.go