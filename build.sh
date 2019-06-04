env GOOS=windows GOARCH=amd64 go build -o build/OpenBooks_WIN.exe main.go
env GOOS=darwin GOARCH=amd64 go build -o build/OpenBooks_MAC main.go
go build -o build/OpenBooks main.go
