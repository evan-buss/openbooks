package mock

import (
	"io/ioutil"
	"log"
	"net"
)

type Config struct {
	Directory   string
	ResultsFile string
}

func Start(config Config) {
	server, err := net.Listen("tcp", ":8080")
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
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Request: %s\n", string(request))

	log.Println("Handling connection")
	conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
	conn.Write([]byte("Content-Type: text/html\r\n"))
	conn.Write([]byte("Connection: close\r\n"))
	conn.Write([]byte("\r\n"))
	conn.Write([]byte("Hello world\r\n"))
	log.Println("Done sending?")
	conn.Close()
}
