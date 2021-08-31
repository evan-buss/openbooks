package mock

import (
	"bufio"
	"log"
	"net"
)

type Config struct {
	Directory   string
	ResultsFile string
}

func Start(config Config) {
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
	log.Printf("Connection received from %s", conn.RemoteAddr().String())
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		request := scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Println(err)
		}

		log.Printf("Request Received: %s\n", request)
	}

	log.Println("Connection closed.")
}
