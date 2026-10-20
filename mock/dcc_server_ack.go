package mock

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
	"time"
)

// AckDccServer emulates a DCC sender that gates its send window on the
// receiver's acknowledgements, which is what most IRC clients do. Without
// acknowledgements the transfer stalls permanently once WindowSize bytes are in
// flight, reproducing the real world 2 MiB stall.
type AckDccServer struct {
	Port       string
	Data       []byte
	WindowSize int
	log        *log.Logger
}

func (dcc *AckDccServer) Start(ready chan<- struct{}) {
	dcc.log = log.New(os.Stdout, "MOCK ACK-DCC: ", 0)

	server, err := net.Listen("tcp", dcc.Port)
	if err != nil {
		panic(err)
	}
	dcc.log.Println("Listening on " + dcc.Port)
	ready <- struct{}{}

	for {
		conn, err := server.Accept()
		if err != nil {
			return
		}
		go dcc.handler(conn)
	}
}

func (dcc *AckDccServer) handler(conn net.Conn) {
	defer conn.Close()

	windowSize := dcc.WindowSize
	if windowSize <= 0 {
		windowSize = 1 << 20
	}

	acked := 0
	sent := 0
	buf := make([]byte, 4096)
	ack := make([]byte, 4)

	for sent < len(dcc.Data) {
		limit := acked + windowSize
		for sent < len(dcc.Data) && sent < limit {
			end := sent + len(buf)
			if end > len(dcc.Data) {
				end = len(dcc.Data)
			}
			n, err := conn.Write(dcc.Data[sent:end])
			if err != nil {
				dcc.log.Println(err)
				return
			}
			sent += n
		}

		if sent >= len(dcc.Data) {
			break
		}

		// The window is full; wait for the receiver to acknowledge before
		// sending more.
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		if _, err := io.ReadFull(conn, ack); err != nil {
			dcc.log.Printf("stalled: window full after %d bytes with no ack: %v", sent, err)
			return
		}
		acked = int(binary.BigEndian.Uint32(ack))
	}

	// Drain outstanding acknowledgements before closing. A TCP socket closed
	// with unread data is reset (RST), which would discard bytes still buffered
	// at the receiver. Real clients read acknowledgements as they arrive.
	conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	for {
		if _, err := conn.Read(buf); err != nil {
			break
		}
	}

	dcc.log.Printf("Done copying %d bytes\n", sent)
}
