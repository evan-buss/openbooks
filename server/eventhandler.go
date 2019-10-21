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

// DownloadSearchResults downloads the search results,
// constructs a JSON response, and sends it over the websocket
func (h Handler) DownloadSearchResults(text string) {
	searchDownloaded := make(chan string)
	// Download the file and wait until it is completed
	go dcc.NewDownload(text, false, false, searchDownloaded)
	// Retrieve the file's location
	fileLocation := <-searchDownloaded
	fmt.Println(fileLocation)
	WS.WriteJSON(SearchResponse{
		MessageType: SEARCH,
		Books:       ParseSearchFile(fileLocation),
	})
}

// DownloadBookFile downloads the search results, constructs
// as JSON response, and sends it over the websocket
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
		log.Println(err)
	}
}

// NoResults is called when the user searches for something that
// is not available. The server sends an appropriate message to client
func (h Handler) NoResults() {
	WS.WriteJSON(IrcErrorResponse{
		MessageType: IRCERROR,
		Status:      "No results found for the query.",
	})
}

// BadServer is called when the user tries to download a file from a
// server that is not available. The server sends an appropriate message
func (h Handler) BadServer() {
	WS.WriteJSON(IrcErrorResponse{
		MessageType: IRCERROR,
		Status:      "Server is not available. Try another one.",
	})
}

// SearchAccepted is called when the search has been accepted but the user
// must wait in the queue for the search to be executed. Send client status
//  update
func (h Handler) SearchAccepted() {
	WS.WriteJSON(WaitResponse{
		MessageType: WAIT,
		Status:      "Search Accepted into the Queue",
	})
}

// MatchesFound is called when the search returns the number of results
// found. Server sends the client a status update
func (h Handler) MatchesFound(num string) {
	WS.WriteJSON(WaitResponse{
		MessageType: WAIT,
		Status:      "Found " + num + " results for your search",
	})
}
