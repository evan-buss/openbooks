#!/bin/bash

echo "Building React App."
cd server/app
npm install
npm run build
cd ../..

echo "Converting React app to binary file."
go get github.com/rakyll/statik
statik -src server/app/build -dest server/
go get -d -v ./...
go install -v ./...

echo "Building binaries for various platforms.";
env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/openbooks.exe
env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/openbooks_mac
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/openbooks_linux
env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o build/openbooks_linux_arm