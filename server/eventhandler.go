package server

import (
	"fmt"

	"github.com/evan-buss/openbooks/dcc"
)

// Handler is the server implementation of the EventHandler interface.
type Handler struct{}

// DownloadSearchResults downloads the search results,
// constructs a JSON response, and sends it over the websocket
func (h Handler) DownloadSearchResults(text string) {
	fmt.Println("DOWNLOAD SEARCH RESULTS HANDLER")
	searchDownloaded := make(chan string)
	// Download the file and wait until it is completed
	go dcc.NewDownload(text, false, searchDownloaded)
	// Retrieve the file's location
	fileLocation := <-searchDownloaded
	// TODO: Process this file and parse it before generating the search response
	WS.WriteJSON(SearchResponse{
		MessageType: SEARCH,
		Books: []BookDetail{
			BookDetail{
				Server: "Oatmeal",
				Author: fileLocation,
				Title:  "Slaughterhouse Five",
				Format: "PDF",
				Size:   "220.0KB",
				Full:   "THIS IS THE FULL DOWNLOAD STRING",
			},
			BookDetail{
				Server: "Pondering42",
				Author: "Kurt Vonnegut",
				Title:  "Slaughterhouse Five",
				Format: "rar",
				Size:   "245.9KB",
				Full:   "THIS IS THE FULL DOWNLOAD STRING",
			},
		},
	})
}

// DownloadBookFile downloads the search results, constructs
// as JSON response, and sends it over the websocket
func (h Handler) DownloadBookFile(text string) {}

// NoResults is called when the user searches for something that
// is not available. The server sends an appropriate message to client
func (h Handler) NoResults() {}

// BadServer is called when the user tries to download a file from a
// server that is not available. The server sends an appropriate message
func (h Handler) BadServer() {}

// SearchAccepted is called when the search has been accepted but the user
// must wait in the queue for the search to be executed. Send client status
//  update
func (h Handler) SearchAccepted() {}

// MatchesFound is called when the search returns the number of results
// found. Server sends the client a status update
func (h Handler) MatchesFound(num string) {}
