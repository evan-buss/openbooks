package dcc

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"regexp"
	"strconv"
)

const delimiterChar = string(byte(0x01))

// There are two types of DCC strings this program accepts.
// Download contains all the necessary DCC info parsed from the DCC SEND string

var (
	ErrInvalidDCCString = errors.New("invalid dcc send string")
	ErrInvalidIP        = errors.New("unable to convert int IP to string")
	ErrMissingBytes     = errors.New("download size didn't match dcc file size. data could be missing")
)

var dccRegex = regexp.MustCompile(`DCC SEND "?(.+[^"])"?\s(\d+)\s+(\d+)\s+(\d+)\s*`)

type Download struct {
	Filename string
	IP       net.IP
	Port     int
	Size     int64
}

func New(filename string, ip net.IP, port int, size int64) *Download {
	return &Download{
		filename,
		ip,
		port,
		size,
	}
}

// ParseString parses the important data of a DCC SEND string
func ParseString(text string) (*Download, error) {
	groups := dccRegex.FindStringSubmatch(text)

	if len(groups) == 0 {
		return nil, ErrInvalidDCCString
	}

	ip, err := string2ip(groups[2])
	if err != nil {
		return nil, err
	}

	size, err := strconv.ParseInt(groups[4], 10, 64)
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(groups[3])
	if err != nil {
		return nil, err
	}

	return &Download{
		Filename: groups[1],
		IP:       ip.To4(),
		Port:     port,
		Size:     size,
	}, nil
}

// Download writes the data contained in the DCC Download
// See https://en.wikipedia.org/wiki/Direct_Client-to-Client#DCC_SEND
func (download Download) Download(writer io.Writer) error {
	// TODO: Maybe specify deadline?
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", download.IP, download.Port))
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

	var received int64
	bytes := make([]byte, 4096)
	for received < download.Size {
		n, err := conn.Read(bytes)
		if err != nil {
			return err
		}

		_, err = writer.Write(bytes[:n])
		if err != nil {
			return err
		}
		received += int64(n)
	}

	if received != download.Size {
		return ErrMissingBytes
	}

	return nil
}

func (download Download) String() string {
	return fmt.Sprintf(`%sDCC SEND "%s" %d %d %d`, delimiterChar, download.Filename, ip2int(download.IP.To4()), download.Port, download.Size)
}

// Convert a given 32 bit IP integer to an IP string
// Ex) 3232235777 -> 192.168.1.1
func string2ip(nn string) (net.IP, error) {
	temp, err := strconv.ParseUint(nn, 10, 32)
	if err != nil {
		return nil, ErrInvalidIP
	}
	intIP := uint32(temp)

	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, intIP)
	return ip, nil
}

func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		panic("no sane way to convert ipv6 into uint32")
	}
	return binary.BigEndian.Uint32(ip)
}
