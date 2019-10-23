packr2
env GOOS=windows GOARCH=amd64 go build -o build/OpenBooks_WIN.exe
env GOOS=darwin GOARCH=amd64 go build -o build/OpenBooks_MAC
env GOOS=linux GOARCH=amd64 go build -o build/OpenBooks_LINUX
packr2 clean
