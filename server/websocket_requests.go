package server

import (
	"encoding/json"
	"fmt"
	"math"
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

	switch message.RequestType {
	case SEARCH:
		obj = new(SearchRequest)
	case DOWNLOAD:
		obj = new(DownloadRequest)
	}

	err := json.Unmarshal(message.Payload, &obj)
	if err != nil {
		server.log.Println("Invalid request payload.")
		c.send <- ErrorResponse{
			Error:   message.RequestType,
			Details: err.Error(),
		}
	}

	switch message.RequestType {
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
	core.Join(c.irc, server.config.Server)
	handler := server.NewIrcEventHandler(c)

	if server.config.Log {
		logger, _, err := util.CreateLogFile(c.irc.Username, server.config.DownloadDir)
		if err != nil {
			server.log.Println(err)
		}
		handler[core.Message] = func(text string) { logger.Println(text) }
	}

	go core.StartReader(c.ctx, c.irc, handler)

	c.send <- ConnectionResponse{
		MessageType: CONNECT,
		Status:      "IRC Server Ready",
		Name:        c.irc.Username,
		Wait:        0,
	}
}

// handle SearchRequests and send the query to the book server
func (c *Client) sendSearchRequest(s *SearchRequest, server *server) {
	server.lastSearchMutex.Lock()
	defer server.lastSearchMutex.Unlock()

	nextAvailableSearch := server.lastSearch.Add(server.config.SearchTimeout)

	if time.Now().Before(nextAvailableSearch) {
		remainingSeconds := time.Until(nextAvailableSearch).Seconds()

		c.send <- SearchRateLimitResponse{
			MessageType: SEARCHRATELIMIT,
			Status:      fmt.Sprintf("Please wait %v seconds to submit another search.", math.Round(remainingSeconds)),
		}
		return
	}

	core.SearchBook(c.irc, s.Query)
	server.lastSearch = time.Now()

	c.send <- WaitResponse{
		MessageType: WAIT,
		Status:      "Search request sent",
	}
}

// handle DownloadRequests by sending the request to the book server
func (c *Client) sendDownloadRequest(d *DownloadRequest) {
	core.DownloadBook(c.irc, d.Book)

	c.send <- WaitResponse{
		MessageType: WAIT,
		Status:      "Download request received",
	}
}
