package dcc

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"regexp"
	"strconv"
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

// downloadDCC downloads the data contained in the Data object
func (download Download) Download(writer io.Writer) error {
	conn, err := net.Dial("tcp", download.IP+":"+download.Port)
	if err != nil {
		return err
	}
	defer conn.Close()

	written, err := io.Copy(writer, conn)
	if err != nil {
		return err
	}

	if written != download.Size {
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
