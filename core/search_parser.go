package core

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"
)

var fileTypes = [...]string{
	"epub",
	"mobi",
	"azw3",
	"html",
	"rtf",
	"pdf",
	"cdr",
	"rar",
	"zip",
}

// BookDetail contains the details of a single Book found on the IRC server
type BookDetail struct {
	Server string `json:"server"`
	Author string `json:"author"`
	Title  string `json:"title"`
	Format string `json:"format"`
	Size   string `json:"size"`
	Full   string `json:"full"`
}

// ParseSearchFile converts a single search file into an array of BookDetail
func ParseSearchFile(filePath string) []BookDetail {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	var books []BookDetail

	counter := 0
	for scanner.Scan() {
		if counter > 3 {
			dat, err := parseLine(scanner.Text())
			if err == nil {
				books = append(books, dat)
			}
		}
		counter++
	}

	//for _, book := range books {
	//	fmt.Printf("%+v\n", book)
	//}

	return books
}

// Parse line extracts data from a single line
func parseLine(line string) (BookDetail, error) {

	//First check if it follows the correct format. Some servers don't include file info...
	if !strings.Contains(line, "::INFO::") {
		return BookDetail{}, errors.New("invalid line format. No file info ")
	}

	var book BookDetail
	book.Full = line[:strings.Index(line, " ::INFO:: ")]
	var tmp int

	// Get Server
	if tmp = strings.Index(line, " "); tmp == -1 {
		return BookDetail{}, errors.New("could not parse server")
	}
	book.Server = line[1:tmp] // Skip the "!"
	line = line[tmp+1:]

	// Get the Author
	if tmp = strings.Index(line, " - "); tmp == -1 {
		return BookDetail{}, errors.New("could not parse author")
	}
	book.Author = line[:tmp]
	line = line[tmp+len(" - "):]

	// Get the Title
	for _, ext := range fileTypes { //Loop through each possible file extension we've got on record
		tmp = strings.Index(line, "."+ext) // check if it contains our extension
		if tmp == -1 {
			continue
		}
		book.Format = ext
		if ext == "rar" || ext == "zip" { // If the extension is .rar or .zip the actual format is contained in ()
			for _, ext2 := range fileTypes[:len(fileTypes)-2] { // Range over the eBook formats (exclude archives)
				if strings.Contains(line[:tmp], ext2) {
					book.Format = ext2
				}
			}
		}
		book.Title = line[:tmp]
		line = line[tmp+len(ext)+1:]
	}

	if book.Title == "" { // Got through the entire loop without finding a single match
		return BookDetail{}, errors.New("could not parse title")
	}

	// Get the Size
	if tmp = strings.Index(line, "::INFO:: "); tmp == -1 {
		return BookDetail{}, errors.New("could not parse size")
	}
	book.Size = line[tmp+9:]

	return book, nil
}
