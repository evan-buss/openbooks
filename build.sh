#!/bin/bash

# Use this script to create production binaries for each platform
cd server/app
npm run build
statik -src ../app/build
cd ../..
env GOOS=windows GOARCH=amd64 go build -o build/openbooks.exe
env GOOS=darwin GOARCH=amd64 go build -o build/openbooks_mac
env GOOS=linux GOARCH=amd64 go build -o build/openbooks_linux
