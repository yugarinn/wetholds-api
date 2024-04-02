#!/usr/bin/env sh

go install github.com/githubnemo/CompileDaemon@latest
DATA_PATH=./data CompileDaemon -log-prefix=false -build "go build -o bin/wetholds-api ./main.go" -command "./bin/wetholds-api" -exclude-dir=".git"
