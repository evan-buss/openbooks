package main

import (
	"fmt"
)

func (server *server) RefreshJob() {
	fmt.Println("Archive Job Running")
	for _, operator := range server.config.IrcServers {
		fmt.Printf("operator: %v\n", operator)
	}
}
