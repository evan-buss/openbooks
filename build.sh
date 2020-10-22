#!/bin/bash

# Use this script to create production binaries for each platform
cd server/app
npm run build
cd ../..
packr2
env GOOS=windows GOARCH=amd64 go build -o build/openbooks.exe
env GOOS=darwin GOARCH=amd64 go build -o build/openbooks_mac
env GOOS=linux GOARCH=amd64 go build -o build/openbooks_linux
env GOOS=linux GOARCH=arm64 go build -o build/openbooks_linuxPI
packr2 clean
