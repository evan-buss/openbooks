package core

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

// List of common file extensions
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

type ParseError struct {
	Line  string
	Error error
}

func (p ParseError) String() string {
	return fmt.Sprintf("Error: %s. Line: %s.", p.Error, p.Line)
}

// ParseSearchFile converts a single search file into an array of BookDetail
func ParseSearchFile(filePath string) ([]BookDetail, []ParseError) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	return ParseSearch(file)
}

func ParseSearch(reader io.Reader) ([]BookDetail, []ParseError) {
	var books []BookDetail
	var errors []ParseError

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "!") {
			dat, err := parseLine(line)
			if err != nil {
				errors = append(errors, ParseError{Line: line, Error: err})
			} else {
				books = append(books, dat)
			}
		}
	}

	sort.Slice(books, func(i, j int) bool { return books[i].Server < books[j].Server })

	return books, errors
}

// Parse line extracts data from a single line
func parseLine(line string) (BookDetail, error) {

	//First check if it follows the correct format. Some servers don't include file info...
	if !strings.Contains(line, "::INFO::") {
		return BookDetail{}, errors.New("invalid line format. ::INFO:: not found")
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

	line = strings.TrimSpace(line)
	splits := strings.Split(line, " ")

	if len(splits) >= 2 {
		book.Size = splits[1]
	}

	return book, nil
}

func ParseSearchV2(reader io.Reader) ([]BookDetail, []ParseError) {
	var books []BookDetail
	var errors []ParseError

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "!") {
			dat, err := parseLineV2(line)
			if err != nil {
				errors = append(errors, ParseError{Line: line, Error: err})
			} else {
				books = append(books, dat)
			}
		}
	}

	sort.Slice(books, func(i, j int) bool { return books[i].Server < books[j].Server })

	return books, errors
}

func parseLineV2(line string) (BookDetail, error) {
	getServer := func(line string) (string, error) {
		if line[0] != '!' {
			return "", errors.New("result lines must start with '!'")
		}

		firstSpace := strings.Index(line, " ")
		if firstSpace == -1 {
			return "", errors.New("unable parse server name")
		}

		return line[1:firstSpace], nil
	}

	getAuthor := func(line string) (string, error) {
		firstSpace := strings.Index(line, " ")
		dashChar := strings.Index(line, " - ")
		if dashChar == -1 {
			return "", errors.New("unable to parse author")
		}

		return line[firstSpace:dashChar], nil
	}

	getTitle := func(line string) (title string, format string, endIndex int) {
		format = ""
		// Get the Title
		for _, ext := range fileTypes { //Loop through each possible file extension we've got on record
			endTitle := strings.Index(line, "."+ext) // check if it contains our extension
			if endTitle == -1 {
				continue
			}
			format = ext
			if ext == "rar" || ext == "zip" { // If the extension is .rar or .zip the actual format is contained in ()
				for _, ext2 := range fileTypes[:len(fileTypes)-2] { // Range over the eBook formats (exclude archives)
					if strings.Contains(line[:endTitle], ext2) {
						format = ext2
					}
				}
			}
			startIndex := strings.Index(line, " - ")
			title = line[startIndex:endTitle]
			endIndex = endTitle
			// line = line[tmp+len(ext)+1:]
		}

		if title == "" { // Got through the entire loop without finding a single match
			return "", "", -1
		}
		return
	}

	getSize := func(line string) (string, int) {
		infoIndex := strings.LastIndex(line, " ::INFO:: ")

		if infoIndex != -1 {
			// TODO: Figure out correct index math. Mind too sleepy....
			endIndex := strings.Index(line[infoIndex:], " ") + len(line) - infoIndex
			size := line[infoIndex:endIndex]
			return size, endIndex
		}

		return "Unknown", len(line)
	}

	server, err := getServer(line)
	if err != nil {
		return BookDetail{}, err
	}

	author, err := getAuthor(line)
	if err != nil {
		return BookDetail{}, err
	}

	title, format, titleIndex := getTitle(line)
	if titleIndex == -1 {
		return BookDetail{}, errors.New("unable to parse title")
	}

	size, endIndex := getSize(line)

	return BookDetail{
		Server: server,
		Author: author,
		Title:  title,
		Format: format,
		Size:   size,
		Full:   line[:endIndex],
	}, nil
}
