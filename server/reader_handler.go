package server

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/evan-buss/openbooks/core"
	"github.com/evan-buss/openbooks/dcc"
)

// Handler is the server implementation of the ReaderHandler interface.
type Handler struct{}

// DownloadSearchResults downloads from DCC server, parses data, and sends data to client
func (h Handler) DownloadSearchResults(text string) {
	searchDownloaded := make(chan string)
	// Download the file and wait until it is completed
	go dcc.NewDownload(text, false, searchDownloaded)
	// Retrieve the file's location
	fileLocation := <-searchDownloaded
	writeJSON(SearchResponse{
		MessageType: SEARCH,
		Books:       core.ParseSearchFile(fileLocation),
	})

	err := os.Remove(fileLocation)
	if err != nil {
		log.Println("Couldn't remove search file", err)
	}
}

// DownloadBookFile downloads the book file and sends it over the websocket
func (h Handler) DownloadBookFile(text string) {
	bookDownloaded := make(chan string)
	go dcc.NewDownload(text, false, bookDownloaded)
	fileLocation := <-bookDownloaded
	fileName := filepath.Base(fileLocation)

	data, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		log.Println("Error reading data from " + fileLocation)
	}

	writeJSON(DownloadResponse{
		MessageType: DOWNLOAD,
		Name:        fileName,
		File:        data,
	})

	err = os.Remove(fileLocation)
	if err != nil {
		log.Println("Couldn't remove book file", err)
	}
}

// NoResults is called when the server returns that nothing was found for the query
func (h Handler) NoResults() {
	writeJSON(IrcErrorResponse{
		MessageType: IRCERROR,
		Status:      "No results found for the query.",
	})
}

// BadServer is called when the requested download fails because the server is not available
func (h Handler) BadServer() {
	writeJSON(IrcErrorResponse{
		MessageType: IRCERROR,
		Status:      "Server is not available. Try another one.",
	})
}

// SearchAccepted is called when the user's query is accepted into the search queue
func (h Handler) SearchAccepted() {
	writeJSON(WaitResponse{
		MessageType: WAIT,
		Status:      "Search accepted into the queue",
	})
}

// MatchesFound is called when the server finds matches for the user's query
func (h Handler) MatchesFound(num string) {
	writeJSON(WaitResponse{
		MessageType: WAIT,
		Status:      "Found " + num + " results for your search",
	})
}
