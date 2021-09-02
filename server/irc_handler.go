package server

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/evan-buss/openbooks/core"
	"github.com/evan-buss/openbooks/dcc"
)

// Client implements the ReaderHandler interface
type IrcHandler struct {
	*Client
	config *Config
	repo   *Repository
}

// DownloadSearchResults downloads from DCC server, parses data, and sends data to client
func (c *IrcHandler) DownloadSearchResults(text string) {
	searchDownloaded := make(chan string)
	// Download the file and wait until it is completed
	go dcc.NewDownload(text, c.config.DownloadDir, searchDownloaded)
	// Retrieve the file's location
	fileLocation := <-searchDownloaded

	log.Printf("%s: Sending search results.\n", c.uuid.String())
	c.send <- SearchResponse{
		MessageType: SEARCH,
		Books:       core.ParseSearchFile(fileLocation),
	}

	err := os.Remove(fileLocation)
	if err != nil {
		log.Printf("%s: Error deleting search results file: %v.\n", c.uuid, err)
	}
}

// DownloadBookFile downloads the book file and sends it over the websocket
func (c *IrcHandler) DownloadBookFile(text string) {
	bookDownloaded := make(chan string)
	go dcc.NewDownload(text, c.config.DownloadDir, bookDownloaded)
	fileLocation := <-bookDownloaded
	fileName := filepath.Base(fileLocation)

	data, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		log.Printf("%s: Error reading book file: %v.\n", c.uuid, err)
	}

	log.Printf("%s: Sending book entitled %s.\n", c.uuid.String(), fileName)
	c.send <- DownloadResponse{
		MessageType: DOWNLOAD,
		Name:        fileName,
		File:        data,
	}

	if !c.config.Persist {
		err = os.Remove(fileLocation)
		if err != nil {
			log.Printf("%s: Error deleting search results file %v.\n", c.uuid, err)
		}
	}
}

// NoResults is called when the server returns that nothing was found for the query
func (c *IrcHandler) NoResults() {
	c.send <- IrcErrorResponse{
		MessageType: IRCERROR,
		Status:      "No results found for the query.",
	}
}

// BadServer is called when the requested download fails because the server is not available
func (c *IrcHandler) BadServer() {
	c.send <- IrcErrorResponse{
		MessageType: IRCERROR,
		Status:      "Server is not available. Try another one.",
	}
}

// SearchAccepted is called when the user's query is accepted into the search queue
func (c *IrcHandler) SearchAccepted() {
	c.send <- WaitResponse{
		MessageType: WAIT,
		Status:      "Search accepted into the queue.",
	}
}

// MatchesFound is called when the server finds matches for the user's query
func (c *IrcHandler) MatchesFound(num string) {
	c.send <- WaitResponse{
		MessageType: WAIT,
		Status:      "Found " + num + " results for your query.",
	}
}

func (c *IrcHandler) PING(url string) {
	c.irc.PONG("irc.irchighway.net")
}

func (c *IrcHandler) ServerList(servers core.IrcServers) {
	c.repo.servers = servers
}
