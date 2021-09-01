package mock

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var logger *log.Logger

type Config struct {
	Directory   string
	ResultsFile string
}

func Start(config Config) {
	logger = log.New(os.Stdout, "MOCK SERVER: ", 0)

	logger.Println("Listening on port 6667")
	server, err := net.Listen("tcp", ":6667")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := server.Accept()
		if err != nil {
			panic(err)
		}
		go handler(conn, config)
	}
}

func handler(conn net.Conn, config Config) {
	logger.Printf("Connection received from %s", conn.RemoteAddr().String())
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		request := scanner.Text()
		if err := scanner.Err(); err != nil {
			logger.Println(err)
		}

		if strings.Contains(request, "@search") {
			fmt.Println("sending search results")
			conn.Write([]byte("NOTICE returned 28 matches"))
		}

		logger.Printf("Request Received: %s\n", request)
	}

	logger.Println("Connection closed.")
}
