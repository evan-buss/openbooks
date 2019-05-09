package dcc

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
)

// Conn rerepesents a DCC connection to a server
type Conn struct {
	dcc      net.Conn
	filename string
	ip       string
	port     string
	size     int
}

// NewDownload parse the string and downloads the given Search results file
func NewDownload(text string, isBook bool) {
	dcc := Conn{}

	// Parse DCC string for important bits
	if isBook {
		dcc.filename, dcc.ip, dcc.port, dcc.size = parseBook(text)
	} else {
		dcc.filename, dcc.ip, dcc.port, dcc.size = parseSearch(text)
	}

	// Convert IP String to IP object to string IP
	intIP, _ := strconv.ParseUint(dcc.ip, 10, 32)
	temp := int2ip(uint32(intIP))
	dcc.ip = temp.String()

	// Establish connection with server
	conn, err := net.Dial("tcp", dcc.ip+":"+dcc.port)
	if err != nil {
		log.Fatal("DCC SEARCH Connection ERROR", err)
	}
	dcc.dcc = conn

	// Create New File based on parsed filename
	f, err := os.Create(dcc.filename)
	if err != nil {
		log.Fatal("Error Creating File", err)
	}
	defer f.Close()

	received := 0
	bytes := make([]byte, 1024)
	for received < dcc.size {
		n, err := dcc.dcc.Read(bytes)
		if err != nil {
			log.Fatal("Error Downloading Data", err)
		}
		f.Write(bytes[:n])
		received += n
	}
	fmt.Println("File downloaded :)")
}

func parseSearch(text string) (filename, ip, port string, size int) {
	re := regexp.MustCompile(`:.+\s+(.+)\s+(\d+)\s+(\d+)\s+(\d+)\s*`)
	groups := re.FindStringSubmatch(text)
	size, _ = strconv.Atoi(groups[4])
	return groups[1], groups[2], groups[3], size
}

func parseBook(text string) (filename, ip, port string, size int) {
	re := regexp.MustCompile(`"(.+)"\s+(\d+)\s+(\d+)\s+(\d+)\s*`)
	groups := re.FindStringSubmatch(text)
	print(groups)
	size, _ = strconv.Atoi(groups[4])
	return groups[1], groups[2], groups[3], size
}

func int2ip(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}
