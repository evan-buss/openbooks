package server

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/evan-buss/openbooks/core"
	"github.com/evan-buss/openbooks/dcc"
	"github.com/evan-buss/openbooks/util"
)

// Client implements the ReaderHandler interface
type IrcHandler struct {
	*Client
	config *Config
	repo   *Repository
}

// DownloadSearchResults downloads from DCC server, parses data, and sends data to client
func (c *IrcHandler) DownloadSearchResults(text string) {
	// Download the file and wait until it is completed
	download, err := dcc.ParseString(text)
	if err != nil {
		c.log.Println(err)
		return
	}

	// Create a new file based on the DCC file name
	dccPath := path.Clean(path.Join(c.config.DownloadDir, download.Filename))
	file, err := os.Create(dccPath)
	if err != nil {
		c.log.Println(err)
		return
	}

	// Download DCC data to the file
	err = download.Download(file)
	if err != nil {
		c.log.Println(err)
		return
	}
	file.Close()

	extractedPath, err := util.ExtractArchive(dccPath)
	c.log.Printf("New Path: %s\n", extractedPath)
	if err != nil {
		c.log.Println(err)
		return
	}

	results, errors := core.ParseSearchFile(extractedPath)
	// Output all errors so parser can be improved over time
	if len(errors) > 0 {
		c.log.Printf("%d Search Result Parsing Errors\n", len(errors))
		for _, err := range errors {
			c.log.Println(err)
		}
	}

	if len(results) == 0 {
		c.NoResults()
		return
	}

	c.log.Printf("Sending %d search results.\n", len(results))
	c.send <- SearchResponse{
		MessageType: SEARCH,
		Books:       results,
	}

	err = os.Remove(extractedPath)
	if err != nil {
		c.log.Printf("Error deleting search results file: %v", err)
	}
}

// DownloadBookFile downloads the book file and sends it over the websocket
func (c *IrcHandler) DownloadBookFile(text string) {
	download, err := dcc.ParseString(text)
	if err != nil {
		c.log.Println(err)
		return
	}

	file, err := os.Create(path.Join(c.config.DownloadDir, download.Filename))
	if err != nil {
		c.log.Println(err)
		return
	}

	filePath, err := filepath.Abs(file.Name())
	if err != nil {
		c.log.Println(err)
		return
	}

	err = download.Download(file)
	if err != nil {
		c.log.Println(err)
		return
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		c.log.Printf("%s: Error reading book file: %v.\n", c.uuid, err)
	}

	c.log.Printf("Sending book entitled '%s'.\n", download.Filename)
	c.send <- DownloadResponse{
		MessageType: DOWNLOAD,
		Name:        download.Filename,
		File:        data,
	}

	if !c.config.Persist {
		err = os.Remove(filePath)
		if err != nil {
			c.log.Printf("Error deleting search results file %s.\n", err)
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

func (c *IrcHandler) Ping() {
	c.irc.Pong(c.config.Server)
}

func (c *IrcHandler) ServerList(servers core.IrcServers) {
	c.repo.servers = servers
}
