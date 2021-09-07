package cli

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/evan-buss/openbooks/core"
	"github.com/evan-buss/openbooks/dcc"
	"github.com/evan-buss/openbooks/irc"
	"github.com/evan-buss/openbooks/util"
)

// Handler is the CLI implementation of the EventHandler interface.
type Handler struct {
	irc         *irc.Conn
	downloadDir string
}

// DownloadSearchResults downloads the search results
// and sends user a response message
func (h Handler) DownloadSearchResults(text string) {
	download, err := dcc.ParseString(text)
	if err != nil {
		log.Println(err)
		return
	}

	filePath := path.Join(h.downloadDir, download.Filename)
	file, err := os.Create(filePath)
	if err != nil {
		log.Println(err)
		return
	}

	err = download.Download(file)
	if err != nil {
		log.Println(err)
		return
	}
	file.Close()

	extractedPath, err := util.ExtractArchive(filePath)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Results location: " + extractedPath)
	menu()
}

// DownloadBookFile downloads the search results and sends
// a user a response message
func (h Handler) DownloadBookFile(text string) {
	download, err := dcc.ParseString(text)
	if err != nil {
		log.Println(err)
		return
	}

	filePath := path.Join(h.downloadDir, download.Filename)
	file, err := os.Create(filePath)
	if err != nil {
		log.Println(err)
		return
	}

	err = download.Download(file)
	if err != nil {
		log.Println(err)
		return
	}
	file.Close()

	extractedPath, err := util.ExtractArchive(filePath)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("File location: " + extractedPath)
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

func (h Handler) Ping() {
	h.irc.Pong("irc.irchighway.net")
}

func (h Handler) ServerList(servers core.IrcServers) {}
