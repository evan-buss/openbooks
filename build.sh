#!/bin/bash

echo "Building React App"
cd server/app
npm install
npm run build
cd ../..

echo "Converting React app to binary file."
go get github.com/rakyll/statik
statik -src server/app/build -dest server/

echo "Building binaries for various platforms";
env GOOS=windows GOARCH=amd64 go build -o build/openbooks.exe
env GOOS=darwin GOARCH=amd64 go build -o build/openbooks_mac
env GOOS=linux GOARCH=amd64 go build -o build/openbooks_linux
env GOOS=linux GOARCH=arm64 go build -o build/openbooks_linux_arm