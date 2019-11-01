#!/bin/bash

# Run this on first repository clone to set up and download required packages
go get -u github.com/gobuffalo/packr/v2/packr2
cd server/app
npm install
npm run integrate