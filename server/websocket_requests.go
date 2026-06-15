package server

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/evan-buss/openbooks/core"
	"github.com/evan-buss/openbooks/util"
)

// RequestHandler defines a generic handle() method that is called when a specific request type is made
type RequestHandler interface {
	handle(c *Client)
}

// messageRouter is used to parse the incoming request and respond appropriately
func (server *server) routeMessage(message Request, c *Client) {
	var obj interface{}

	switch message.MessageType {
	case SEARCH:
		obj = new(SearchRequest)
	case DOWNLOAD:
		obj = new(DownloadRequest)
	}

	err := json.Unmarshal(message.Payload, &obj)
	if err != nil {
		server.log.Printf("Invalid request payload. %s.\n", err.Error())
		c.send <- StatusResponse{
			MessageType:      STATUS,
			NotificationType: DANGER,
			Title:            "Unknown request payload.",
		}
	}

	switch message.MessageType {
	case CONNECT:
		c.startIrcConnection(server)
	case SEARCH:
		c.sendSearchRequest(obj.(*SearchRequest), server)
	case DOWNLOAD:
		c.sendDownloadRequest(obj.(*DownloadRequest))
	default:
		server.log.Println("Unknown request type received.")
	}
}

// handle ConnectionRequests and either connect to the server or do nothing
func (c *Client) startIrcConnection(server *server) {
	err := c.bookClient.Connect(server.config.Server, server.config.EnableTLS)
	if err != nil {
		c.log.Println(err)
		c.send <- newErrorResponse("Unable to connect to IRC server.")
		return
	}

	handlers := core.Handlers{
		SearchResults: func(books []core.BookDetail, parseErrors []core.ParseError, resultsFilePath string) {
			if len(parseErrors) > 0 {
				c.log.Printf("%d Search Result Parsing Errors\n", len(parseErrors))
			}
			c.log.Printf("Sending %d search results.\n", len(books))
			c.send <- newSearchResponse(books, parseErrors)
			if err := os.Remove(resultsFilePath); err != nil {
				c.log.Printf("Error deleting search results file: %v", err)
			}
		},
		BookDownloaded: func(filePath string) {
			c.log.Printf("Sending book entitled '%s'.\n", filepath.Base(filePath))
			c.send <- newDownloadResponse(filePath, server.config.DisableBrowserDownloads)
		},
		NoResults:      func() { c.send <- newErrorResponse("No results found for the query.") },
		BadServer:      func() { c.send <- newErrorResponse("Server is not available. Try another one.") },
		SearchAccepted: func() { c.send <- newStatusResponse(NOTIFY, "Search accepted into the queue.") },
		MatchesFound: func(num string) {
			c.send <- newStatusResponse(NOTIFY, fmt.Sprintf("Found %s results for your query.", num))
		},
		ServerList: func(servers core.IrcServers) { server.repository.servers = servers },
		Erred: func(err error) {
			c.log.Println(err)
			c.send <- newErrorResponse("Error processing request.")
		},
	}

	if server.config.Log {
		logger, _, err := util.CreateLogFile(c.bookClient.Username(), server.config.DownloadDir)
		if err != nil {
			server.log.Println(err)
		} else {
			handlers.Message = func(text string) { logger.Println(text) }
		}
	}

	c.bookClient.StartReader(c.ctx, handlers)

	c.send <- ConnectionResponse{
		StatusResponse: StatusResponse{
			MessageType:      CONNECT,
			NotificationType: SUCCESS,
			Title:            "Welcome, connection established.",
			Detail:           fmt.Sprintf("IRC username %s", c.bookClient.Username()),
		},
		Name: c.bookClient.Username(),
	}
}

// handle SearchRequests and send the query to the book server
func (c *Client) sendSearchRequest(s *SearchRequest, server *server) {
	server.lastSearchMutex.Lock()
	defer server.lastSearchMutex.Unlock()

	nextAvailableSearch := server.lastSearch.Add(server.config.SearchTimeout)

	if time.Now().Before(nextAvailableSearch) {
		remainingSeconds := time.Until(nextAvailableSearch).Seconds()
		c.send <- newRateLimitResponse(remainingSeconds)

		return
	}

	c.bookClient.Search(server.config.SearchBot, s.Query)
	server.lastSearch = time.Now()

	c.send <- newStatusResponse(NOTIFY, "Search request sent.")
}

// handle DownloadRequests by sending the request to the book server
func (c *Client) sendDownloadRequest(d *DownloadRequest) {
	c.bookClient.Download(d.Book)
	c.send <- newStatusResponse(NOTIFY, "Download request received.")
}
