.ONESHELL:
.PHONY: init mock dev-cli dev-server dev-client build update

init:
	go get -u ./...
	cd server/app
	npm install

mock:
	cd cmd/mock_server
	go build
	./mock_server

dev-cli:
	cd cmd/openbooks
	go build
	./openbooks cli --server localhost

dev-server:
	cd cmd/openbooks
	go build
	./openbooks server --server localhost

dev-api:
	cd cmd/openbooks_api
	go build
	./openbooks_api

dev-client:
	cd server/app
	npm run dev

build:
	rm -r ./build
	./build.sh

update:
	go get -u
	go mod tidy
	cd server/app
	npx npm-check-updates -i