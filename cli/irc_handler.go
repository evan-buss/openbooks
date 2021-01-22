package cli

import (
	"fmt"

	"github.com/evan-buss/openbooks/dcc"
)

// Handler is the CLI implementation of the EventHandler interface.
type Handler struct {
	downloadDir string
}

// DownloadSearchResults downloads the search results
// and sends user a response message
func (h Handler) DownloadSearchResults(text string) {
	fileLocation := make(chan string)
	go dcc.NewDownload(text, h.downloadDir, fileLocation)
	fmt.Println("Results location: " + <-fileLocation)
	menu()
}

// DownloadBookFile downloads the search results and sends
// a user a response message
func (h Handler) DownloadBookFile(text string) {
	fileLocation := make(chan string)
	go dcc.NewDownload(text, h.downloadDir, fileLocation)
	fmt.Println("File location: " + <-fileLocation)
	menu()
}

// NoResults is called when the user searches for something that
// is not available sends a CLI message
func (h Handler) NoResults() {
	fmt.Println("No results returned for that search...")
	menu()
}

// BadServer is called when the user tries to download a file from a
// server that is not available.
func (h Handler) BadServer() {
	fmt.Println("That server is not available. Try again...")
	menu()
}

// SearchAccepted is called when the search has been accepted but the user
// must wait in the queue for the search to be executed.
func (h Handler) SearchAccepted() {
	fmt.Println("Search has been accepted. Please wait.")
}

// MatchesFound is called when the search returns the number of results
// found. Server sends the client a status update
func (h Handler) MatchesFound(num string) {
	fmt.Println("Your search returned " + num + " matches.")
}
