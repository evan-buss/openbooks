package dcc

import (
	"bufio"
	"encoding/binary"
	"errors"
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

// Conn represents a DCC connection to a server
type Conn struct {
	dcc      net.Conn
	filename string
	ip       string
	port     string
	size     int
}

// TODO: Refactor this function. It is kind of messy as is. I should use the filepath utilities instead of manually handling paths
// NewDownload parses the string and downloads the file
func NewDownload(text string, isBook bool, isCli bool, doneChan chan<- string) {
	pathSep := string(os.PathSeparator)
	dcc := Conn{}

	var err error
	downloadDir := os.TempDir() + pathSep

	if isCli {
		downloadDir, err = os.UserHomeDir()

		if err != nil {
			log.Fatal("Home Dir Error: ", err)
		}

		downloadDir += pathSep + "Downloads" + pathSep
	}

	// Parse DCC string for important bits
	if isBook {
		err := dcc.ParseBook(text)
		if err != nil {
			log.Fatal("Parse Books: ", err)
		}
	} else {
		err := dcc.ParseSearch(text)
		dcc.filename = "search_results.zip"
		if err != nil {
			log.Fatal("Parse Books: ", err)
		}
	}

	// Establish connection with server
	conn, err := net.Dial("tcp", dcc.ip+":"+dcc.port)
	if err != nil {
		log.Fatal("DCC Connection ERROR", err)
	}
	dcc.dcc = conn

	// Create New File based on parsed filename
	zipfile, err := os.Create(downloadDir + dcc.filename)
	if err != nil {
		log.Fatal("Error Creating File", err)
	}

	fileExt := filepath.Ext(zipfile.Name())
	if fileExt == ".rar" || fileExt == ".zip" {
		defer os.Remove(zipfile.Name()) // We don't want the zip file so delete when done
	}
	defer zipfile.Close()

	// Download the file
	received := 0
	bytes := make([]byte, 1024)
	for received < dcc.size {
		n, err := dcc.dcc.Read(bytes)
		if err != nil {
			log.Fatal("Error Downloading Data", err)
		}
		zipfile.Write(bytes[:n])
		received += n
	}
	fName := zipfile.Name()
	// For each file in the archive, save the data to a new file (extract)
	err = archiver.Walk(downloadDir+dcc.filename, func(f archiver.File) error {
		fName = f.Name()

		file, err := os.OpenFile(downloadDir+fName, os.O_CREATE|os.O_RDWR, 0644)
		defer file.Close()
		if err != nil {
			log.Println("Error Creating TXT: " + fName)
		}

		writer := bufio.NewWriter(file)
		reader := bufio.NewReader(f)

		for {
			// Read from the archived file and write to the new file name
			data := make([]byte, 1024)
			n, err := reader.Read(data)

			writer.Write(data[:n])
			if err != nil {
				writer.Flush()
				break
			}
		}
		return nil
	})

	// Pass the channel the location of the newly download file
	doneChan <- downloadDir + fName
}

// ParseSearch parses the important data from a search results string
func (dcc *Conn) ParseSearch(text string) error {
	re := regexp.MustCompile(`:.+\s+(.+)\s+(\d+)\s+(\d+)\s+(\d+)\s*`)
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

// ParseBook parses the important data of a book download string
func (dcc *Conn) ParseBook(text string) error {
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
