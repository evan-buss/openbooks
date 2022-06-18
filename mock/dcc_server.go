package mock

import (
	"io"
	"log"
	"net"
	"os"
	"time"
)

type DccServer struct {
	Port   string
	Reader io.ReadSeeker
	log    *log.Logger
}

func (dcc *DccServer) Start(ready chan<- struct{}) {
	dcc.log = log.New(os.Stdout, "MOCK DCC: ", 0)

	server, err := net.Listen("tcp", dcc.Port)
	if err != nil {
		panic(err)
	}
	dcc.log.Println("Listening on " + dcc.Port)
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
		dcc.log.Println("closing connection")
		conn.Close()
	}()

	dcc.log.Println("Received a connection...")

	var err error
	var n int
	// n, err := io.Copy(conn, dcc.Reader)

	// Use below to slow download speed for testing
	bytes := make([]byte, 4096)
	for {
		n, _ = dcc.Reader.Read(bytes)
		_, err = conn.Write(bytes[:n])

		time.Sleep(time.Millisecond * 250)

		if n == 0 {
			break
		}
	}

	if err != nil {
		dcc.log.Println(err)
	} else {
		dcc.log.Printf("Done copying %d bytes\n", n)
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
