cd server/app
npm run build
cd ../..
packr2
env GOOS=windows GOARCH=amd64 go build -o build/openbooks.exe
env GOOS=darwin GOARCH=amd64 go build -o build/openbooks_mac
env GOOS=linux GOARCH=amd64 go build -o build/openbooks_linux
packr2 clean
