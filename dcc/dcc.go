package dcc

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/mholt/archiver"
)

// There are two types of DCC strings this program accepts.
// Search Results
// 	- A text file containing a list of search results returned from your search
//    query.
// Book Files
//  - The actual book file itself. You get the download string from the search
//    results, enter it, then the book is sent to you.

// Data contains all of the necessary DCC info parsed from the DCC SEND string
type Data struct {
	filename string
	ip       string
	port     string
	size     int
}

// NewDownload parses the string and downloads the file
func NewDownload(text string, isBook bool, isCli bool, doneChan chan<- string) {
	dcc := Data{}

	err := dcc.ParseDCC(text)
	if err != nil {
		log.Fatal("ParseDCC Error: ", err)
	}

	var directory string

	if isCli {
		directory, err = os.UserHomeDir()
		if err != nil {
			log.Println(err)
		}
		directory = filepath.Join(directory, "Downloads")
	} else {
		directory = os.TempDir()
	}

	fullPath := filepath.Join(directory, dcc.filename)

	file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.Dial("tcp", dcc.ip+":"+dcc.port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Download the file
	received := 0
	bytes := make([]byte, 1024)
	for received < dcc.size {
		n, err := conn.Read(bytes)
		if err != nil {
			log.Fatal("Error Downloading Data", err)
		}
		_, err = file.Write(bytes[:n])
		if err != nil {
			log.Println(err)
		}
		received += n
	}

	var newPath string
	ext := filepath.Ext(fullPath)
	if ext == ".rar" || ext == ".zip" {
		zipPath := fullPath // Save the zip file path and keep it unchanged
		defer os.Remove(zipPath)

		err = archiver.Walk(fullPath, func(f archiver.File) error {
			newPath = filepath.Join(filepath.Dir(fullPath), f.Name())

			out, err := os.OpenFile(newPath, os.O_CREATE|os.O_RDWR, 0644)
			if err != nil {
				return err
			}

			_, err = io.Copy(out, f)
			if err != nil {
				return err
			}

			err = out.Close()
			if err != nil {
				log.Println("Error closing the archive output file", err)
			}

			return nil
		})
		if err != nil {
			log.Println(err)
		}
	}
	if newPath != "" { //If we extracted a file, update the output path
		fullPath = newPath
	}
	doneChan <- fullPath
}

// ParseDCC parses the important data of a book download string
func (dcc *Data) ParseDCC(text string) error {
	re := regexp.MustCompile(`DCC SEND "?(.+[^"])"?\s(\d+)\s+(\d+)\s+(\d+)\s*`)
	groups := re.FindStringSubmatch(text)

	if len(groups) == 0 {
		return errors.New("no match in string")
	}

	dcc.filename = groups[1]
	dcc.ip = stringToIP(groups[2])
	dcc.port = groups[3]
	dcc.size, _ = strconv.Atoi(groups[4])
	return nil
}

// Convert a given 32 bit IP integer to an IP string
// Ex) 2907707975 -> 192.168.1.1
func stringToIP(nn string) string {
	temp, err := strconv.ParseUint(nn, 10, 32)
	if err != nil {
		log.Println("Error Parsing Int From Host String: ", err)
	}
	intIP := uint32(temp)

	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, intIP)
	return ip.String()
}
