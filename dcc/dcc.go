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

// TODO: Clean this up and use built in io methods

// There are two types of DCC strings this program accepts.
// Search Results
// 	- A text file containing a list of search results returned from a search
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

// NewDownload parses the DCC SEND string and downloads the file
func NewDownload(text string, downloadDir string, doneChan chan<- string) {
	dcc := Data{}
	err := dcc.parseDCC(text)
	if err != nil {
		// TODO: This shouldn't be fatal. Need a good solution for error logging...
		log.Fatal("ParseDCC Error: ", err)
	}

	dcc.filename = filepath.Join(downloadDir, dcc.filename)

	downloadDCC(dcc, doneChan)
}

// downloadDCC downloads the data contained in the Data object
func downloadDCC(dcc Data, doneChan chan<- string) {

	file, err := os.OpenFile(dcc.filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.Dial("tcp", dcc.ip+":"+dcc.port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Download the file
	// NOTE: I tried the io.Copy utility but it EASILY took about 100x
	// the amount of time as the manual method... Not sure why (windows)
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

	file.Close()

	var newPath string
	ext := filepath.Ext(dcc.filename)
	if ext == ".rar" || ext == ".zip" {
		err = archiver.Walk(dcc.filename, func(f archiver.File) error {
			newPath = filepath.Join(filepath.Dir(dcc.filename), f.Name())

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

	if newPath != "" { // If we extracted a file, send that file and remove the zip file
		doneChan <- newPath
		err = os.Remove(dcc.filename)
		if err != nil {
			log.Println("remove error", err)
		}
	} else {
		doneChan <- dcc.filename
	}
}

// ParseDCC parses the important data of a DCC SEND string
func (dcc *Data) parseDCC(text string) error {
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
