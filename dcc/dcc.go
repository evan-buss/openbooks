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

// dialTimeout bounds how long we wait to establish the DCC connection.
const dialTimeout = 30 * time.Second

// idleTimeout is the maximum time to wait for more data during a transfer.
// It is reset after every successful read, so it only fires when the sender
// stalls (stops sending) rather than limiting the total transfer duration.
const idleTimeout = 60 * time.Second

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

// Download writes the data contained in the DCC Download
func (download Download) Download(writer io.Writer) error {
	conn, err := net.DialTimeout("tcp", download.IP+":"+download.Port, dialTimeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	// NOTE: Not using the idiomatic io.Copy or io.CopyBuffer because they are
	// much slower in real world tests than the manual way. I suspect it has to
	// do with the way the DCC server is sending data. I don't think it ever sends
	// an EOF like the io.* methods expect.

	// Benchmark: 2.36MB File
	// CopyBuffer - 4096 - 2m32s, 2m18s, 2m32s
	// Copy - 2m35s
	// Custom - 1024 - 35s
	// Custom - 4096 - 46s, 14s
	received := 0
	bytes := make([]byte, 4096)
	for int64(received) < download.Size {
		// Reset the idle deadline before each read so a stalled sender causes
		// the transfer to fail instead of blocking forever.
		if err := conn.SetReadDeadline(time.Now().Add(idleTimeout)); err != nil {
			return err
		}

		n, err := conn.Read(bytes)
		if err != nil {
			return err
		}

		_, err = writer.Write(bytes[:n])
		if err != nil {
			return err
		}
		received += n
	}

	if int64(received) != download.Size {
		return ErrMissingBytes
	}

	return nil
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
