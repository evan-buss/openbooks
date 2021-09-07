package mock

import (
	"io"
	"log"
	"net"
	"os"
)

type DccServer struct {
	Port   string
	Reader io.ReadSeeker
}

func (dcc *DccServer) Start(ready chan<- struct{}) {
	logger := log.New(os.Stdout, "MOCK DCC: ", 0)

	server, err := net.Listen("tcp", dcc.Port)
	if err != nil {
		panic(err)
	}
	logger.Println("Listening on " + dcc.Port)
	ready <- struct{}{}

	for {
		conn, err := server.Accept()
		if err != nil {
			panic(err)
		}
		go dcc.handler(conn)
	}
}

func (dcc *DccServer) handler(conn net.Conn) {
	defer func() {
		dcc.Reader.Seek(0, io.SeekStart)
		logger.Println("closing connection")
		conn.Close()
	}()

	logger.Println("Received a connection...")
	n, err := io.Copy(conn, dcc.Reader)
	if err != nil {
		logger.Println(err)
	} else {
		logger.Printf("Done copying %d bytes\n", n)
	}
}

type WriteCloser struct {
	Data []byte
}

func (m *WriteCloser) Write(p []byte) (n int, err error) {
	m.Data = append(m.Data, p...)
	return len(p), nil
}

func (m WriteCloser) Close() error {
	return nil
}
