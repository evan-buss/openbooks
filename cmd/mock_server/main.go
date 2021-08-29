package main

import "github.com/evan-buss/openbooks/mock"

func main() {
	config := mock.Config{}
	mock.Start(config)
}
