package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/evan-buss/openbooks/dcc"
)

// Handler is the server implementation of the EventHandler interface.
type Handler struct{}

// DownloadSearchResults downloads from DCC server, parses data, and sends data to client
func (h Handler) DownloadSearchResults(text string) {
	searchDownloaded := make(chan string)
	// Download the file and wait until it is completed
	go dcc.NewDownload(text, false, false, searchDownloaded)
	// Retrieve the file's location
	fileLocation := <-searchDownloaded
	fmt.Println(fileLocation)
	err := WS.WriteJSON(SearchResponse{
		MessageType: SEARCH,
		Books:       ParseSearchFile(fileLocation),
	})
	if err != nil {
		log.Println("Error Sending SearchResponse: ", err)
	}
}

// DownloadBookFile downloads the book file and sends it over the websocket
func (h Handler) DownloadBookFile(text string) {
	bookDownloaded := make(chan string)
	go dcc.NewDownload(text, true, false, bookDownloaded)
	// Wait until the download finishes and get the location
	fileLocation := <-bookDownloaded

	fileName := filepath.Base(fileLocation)

	fmt.Println(fileLocation)
	fmt.Println(fileName)

	data, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		// TODO: Send an error message to the client
		log.Println("Error reading data from " + fileLocation)
	}

	err = WS.WriteJSON(DownloadResponse{
		MessageType: DOWNLOAD,
		Name:        fileName,
		File:        data,
	})
	if err != nil {
		log.Println("Error sending DownloadResponse: ", err)
	}
}

// NoResults is called when the server returns that nothing was found for the query
func (h Handler) NoResults() {
	err := WS.WriteJSON(IrcErrorResponse{
		MessageType: IRCERROR,
		Status:      "No results found for the query.",
	})
	if err != nil {
		log.Println("Error sending IrcErrorResponse: ", err)
	}
}

// BadServer is called when the requested download fails because the server is not available
func (h Handler) BadServer() {
	err := WS.WriteJSON(IrcErrorResponse{
		MessageType: IRCERROR,
		Status:      "Server is not available. Try another one.",
	})
	if err != nil {
		log.Println("Error sending IrcServerResponse: ", err)
	}
}

// SearchAccepted is called when the user's query is accepted into the search queue
func (h Handler) SearchAccepted() {
	err := WS.WriteJSON(WaitResponse{
		MessageType: WAIT,
		Status:      "Search Accepted into the Queue",
	})
	if err != nil {
		log.Println("Error sending WaitResponse: ", err)
	}
}

// MatchesFound is called when the server finds matches for the user's query
func (h Handler) MatchesFound(num string) {
	err := WS.WriteJSON(WaitResponse{
		MessageType: WAIT,
		Status:      "Found " + num + " results for your search",
	})
	if err != nil {
		log.Println("Error sending WaitResponse: ", err)
	}
}
