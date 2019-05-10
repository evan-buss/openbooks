package dcc

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/mholt/archiver"
)

// Default Locations
const searchFolder string = "./Searches/"
const bookFolder string = "./Books/"

// TODO: File downloads and extractions are fucked. Something to do with paths i think

// Conn rerepesents a DCC connection to a server
type Conn struct {
	dcc      net.Conn
	filename string
	ip       string
	port     string
	size     int
}

// NewDownload parses the string and downloads the file
func NewDownload(text string, isBook bool, doneChan chan<- bool) {
	dcc := Conn{}

	defer func() {
		// Send message to continue listening when finished downloading
		doneChan <- true
	}()

	// Parse DCC string for important bits
	if isBook {
		err := dcc.ParseBook(text)
		if err != nil {
			log.Fatal("Parse Books: ", err)
		}
	} else {
		err := dcc.ParseSearch(text)
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
	zipfile, err := os.Create(dcc.filename)
	if err != nil {
		log.Fatal("Error Creating File", err)
	}

	fileExt := filepath.Ext(zipfile.Name())
	if fileExt == ".rar" || fileExt == ".zip" {
		defer os.Remove(zipfile.Name()) // We don't want the zip file so delete when done
	}
	defer zipfile.Close()

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

	// For each file in the archive, save the data to a new file (extract)
	err = archiver.Walk(dcc.filename, func(f archiver.File) error {
		fName := f.Name()

		file, err := os.Create(fName)
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
				fmt.Println(err)
				fmt.Println("BUFFER SIZE: " + strconv.Itoa(writer.Buffered()))
				writer.Flush()
				break
			}
		}
		fmt.Println("File has been downloaded")
		return nil
	})
}

// ParseSearch parses the important data from a search results string
func (dcc *Conn) ParseSearch(text string) error {
	re := regexp.MustCompile(`:.+\s+(.+)\s+(\d+)\s+(\d+)\s+(\d+)\s*`)
	groups := re.FindStringSubmatch(text)

	if len(groups) == 0 {
		return errors.New("No match in string")
	}

	dcc.filename = groups[1]
	dcc.ip = stringToIP(groups[2])
	dcc.port = groups[3]
	dcc.size, _ = strconv.Atoi(groups[4])
	print("Parsed search string")

	return nil
}

// ParseBook parses the important data of a book download string
func (dcc *Conn) ParseBook(text string) error {
	re := regexp.MustCompile(`DCC SEND (.+)\s+(\d+)\s+(\d+)\s+(\d+)\s*`)
	// re := regexp.MustCompile(`"(.+)"\s+(\d+)\s+(\d+)\s+(\d+)\s*`)
	groups := re.FindStringSubmatch(text)

	if len(groups) == 0 {
		return errors.New("No match in string")
	}

	dcc.filename = groups[1]
	dcc.ip = stringToIP(groups[2])
	dcc.port = groups[3]
	dcc.size, _ = strconv.Atoi(groups[4])
	print("Parsed book string")
	return nil
}

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
