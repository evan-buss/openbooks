package dcc

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"regexp"
	"strconv"
	"time"
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
	// TODO: Maybe specify deadline?
	conn, err := net.Dial("tcp", download.IP+":"+download.Port)
	if err != nil {
		return err
	}

	start := time.Now()
	defer func() {
		log.Printf("DCC: %d bytes took %s\n to download.\n", download.Size, time.Since(start))
		conn.Close()
	}()

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
		n, err := conn.Read(bytes)

		if err != nil {
			log.Fatal("Error Downloading Data", err)
		}
		_, err = writer.Write(bytes[:n])
		if err != nil {
			log.Println(err)
		}
		received += n
	}

	if int64(received) != download.Size {
		return ErrMissingBytes
	}

	return nil
}

// ParseDCC parses the important data of a DCC SEND string
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
