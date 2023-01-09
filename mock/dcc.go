package mock

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type DccServer struct {
	file string
	log  *log.Logger
}

func NewDccServer(file string) *DccServer {
	return &DccServer{
		file: file,
		log:  log.New(os.Stdout, "DCC: ", 0),
	}
}

func (dcc *DccServer) Start(selectedPort chan<- int) {
	var (
		server net.Listener
		err    error
		port   = 6667
	)
	for {
		server, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			break
		}
		port++
	}
	defer server.Close()
	selectedPort <- port
	dcc.log = log.New(os.Stdout, fmt.Sprintf("DCC (%d): ", port), 0)

	dcc.log.Printf("Listening on %d\n", port)

	conn, err := server.Accept()

	dcc.log.Println("Received a connection")
	if err != nil {
		panic(err)
	}
	dcc.handler(conn)

	dcc.log.Println("Closing Server")
}

func (dcc *DccServer) handler(conn net.Conn) {
	defer conn.Close()
	dcc.log.Printf("Sending file %s\n", dcc.file)
	file, err := os.Open(dcc.file)
	if err != nil {
		dcc.log.Println(err)
		return
	}
	defer file.Close()

	//confirmation := make([]byte, 4)
	for {
		n, err := io.CopyN(conn, file, 1024)
		//conn.Read(confirmation)

		//time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
		//log.Println(binary.BigEndian.Uint32(confirmation))

		if n == 0 || err != nil {
			break
		}
	}

	if err != nil && !errors.Is(err, io.EOF) {
		dcc.log.Println(err)
		return
	}

	dcc.log.Println("Done sending file.")
}
