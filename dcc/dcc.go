package dcc

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"regexp"
	"strconv"
	"time"
)

const (
	// ackInterval is how many bytes may arrive before the receiver must send a
	// DCC "SENT" acknowledgement. Senders gate their send window on these; a
	// receiver that never sends them stalls the transfer once the window fills.
	ackInterval = 1024

	// stallTimeout aborts a transfer when the sender stops sending for this
	// long. Without it a stalled sender leaves a permanent .temp file behind
	// and the caller never learns the download died.
	stallTimeout = 5 * time.Minute

	// dialTimeout bounds how long we wait for the sender's DCC listener to
	// accept the connection. net.Dial without a timeout can hang indefinitely
	// on an unreachable sender.
	dialTimeout = 30 * time.Second
)

// There are two types of DCC strings this program accepts.
// Download contains all of the necessary DCC info parsed from the DCC SEND string

var (
	ErrInvalidDCCString = errors.New("invalid dcc send string")
	ErrInvalidIP        = errors.New("unable to convert int IP to string")
	ErrMissingBytes     = errors.New("download size didn't match dcc file size. data could be missing")
)

var dccRegex = regexp.MustCompile(`DCC SEND "?(.+[^"])"?\s(\d+)\s+(\d+)\s+(\d+)\s*`)

type Download struct {
	Filename string
	IP       string
	Port     string
	Size     int64
}

// ParseString parses the important data of a DCC SEND string
func ParseString(text string) (*Download, error) {
	groups := dccRegex.FindStringSubmatch(text)

	if len(groups) == 0 {
		return nil, ErrInvalidDCCString
	}

	ip, err := stringToIP(groups[2])
	if err != nil {
		return nil, err
	}

	size, err := strconv.ParseInt(groups[4], 10, 64)
	if err != nil {
		return nil, err
	}

	return &Download{
		Filename: groups[1],
		IP:       ip,
		Port:     groups[3],
		Size:     size,
	}, nil
}

// receiver is the subset of net.Conn the transfer loop depends on. Using an
// interface keeps the protocol logic testable without a real socket.
type receiver interface {
	io.Reader
	io.Writer
	SetReadDeadline(time.Time) error
}

// Download writes the data contained in the DCC Download
func (download Download) Download(writer io.Writer) error {
	conn, err := net.DialTimeout("tcp", download.IP+":"+download.Port, dialTimeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	received, err := receive(conn, writer, download.Size)
	if err != nil {
		return err
	}
	if received != download.Size {
		return ErrMissingBytes
	}

	return nil
}

// receive copies exactly size bytes from conn into writer while acknowledging
// the sender's window. It returns the number of bytes written even when it
// returns an error, so a caller can tell how much actually arrived.
func receive(conn receiver, writer io.Writer, size int64) (int64, error) {
	// NOTE: Not using the idiomatic io.Copy or io.CopyBuffer because they are
	// much slower in real world tests than the manual way. I suspect it has to
	// do with the way the DCC server is sending data. I don't think it ever sends
	// an EOF like the io.* methods expect.

	// Benchmark: 2.36MB File
	// CopyBuffer - 4096 - 2m32s, 2m18s, 2m32s
	// Copy - 2m35s
	// Custom - 1024 - 35s
	// Custom - 4096 - 46s, 14s
	var received, acked int64
	bytes := make([]byte, 4096)
	for received < size {
		// The deadline is refreshed on every read, so slow but steady transfers
		// finish while a sender that goes silent is abandoned.
		if err := conn.SetReadDeadline(time.Now().Add(stallTimeout)); err != nil {
			return received, err
		}

		n, err := conn.Read(bytes)

		// A Read may return data together with an error. Writing the data first
		// is essential: senders that reset the connection right after the final
		// bytes (a common XDCC failure) would otherwise lose that tail and the
		// transfer would be reported as incomplete.
		if n > 0 {
			if _, werr := writer.Write(bytes[:n]); werr != nil {
				return received, werr
			}
			received += int64(n)

			// Open the sender's window by reporting progress. Senders with a
			// limited window (commonly 1-2 MiB) stop permanently otherwise.
			if received-acked >= ackInterval || received == size {
				if aerr := sendAck(conn, received); aerr != nil {
					return received, aerr
				}
				acked = received
			}
		}

		if err != nil {
			// If the final bytes and the error arrived together, the download
			// is actually complete.
			if received == size {
				return received, nil
			}
			return received, err
		}
	}

	return received, nil
}

// sendAck reports the number of bytes received so far, as required by the DCC
// SEND protocol.
func sendAck(w io.Writer, received int64) error {
	var ack [4]byte
	binary.BigEndian.PutUint32(ack[:], uint32(received))
	_, err := w.Write(ack[:])
	return err
}

// Convert a given 32 bit IP integer to an IP string
// Ex) 2907707975 -> 192.168.1.1
func stringToIP(nn string) (string, error) {
	temp, err := strconv.ParseUint(nn, 10, 32)
	if err != nil {
		return "", ErrInvalidIP
	}
	intIP := uint32(temp)

	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, intIP)
	return ip.String(), nil
}
