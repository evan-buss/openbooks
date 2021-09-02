package server

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/evan-buss/openbooks/core"
)

// RequestHandler defines a generic handle() method that is called when a specific request type is made
type RequestHandler interface {
	handle(c *Client)
}

// messageRouter is used to parse the incoming request and respond appropriately
func (s *server) routeMessage(message Request, c *Client) {
	var obj interface{}

	switch message.RequestType {
	case SEARCH:
		obj = new(SearchRequest)
	case DOWNLOAD:
		obj = new(DownloadRequest)
	}

	err := json.Unmarshal(message.Payload, &obj)
	if err != nil {
		log.Printf("%s: Invalid request payload.\n", c.uuid.String())
		c.send <- ErrorResponse{
			Error:   message.RequestType,
			Details: err.Error(),
		}
	}
	fmt.Printf("%#v\n", obj)

	switch message.RequestType {
	case CONNECT:
		c.handleConnectionRequest(s)
	case SEARCH:
		c.handleSearchRequest(obj.(*SearchRequest))
	case DOWNLOAD:
		c.handleDownloadRequest(obj.(*DownloadRequest))
	default:
		log.Printf("%s: Unknown request type received.\n", c.uuid.String())
	}
}

// handle ConnectionRequests and either connect to the server or do nothing
func (c *Client) handleConnectionRequest(server *server) {
	core.Join(c.irc)
	ircHandler := &IrcHandler{c, server.config, server.repository}
	go core.ReadDaemon(c.irc, ircHandler, false, c.disconnect)

	c.send <- ConnectionResponse{
		MessageType: CONNECT,
		Status:      "IRC Server Ready",
		Wait:        0,
	}
}

// handle SearchRequests and send the query to the book server
func (c *Client) handleSearchRequest(s *SearchRequest) {
	core.SearchBook(c.irc, s.Query)

	c.send <- WaitResponse{
		MessageType: WAIT,
		Status:      "Search request sent",
	}
}

// handle DownloadRequests by sending the request to the book server
func (c *Client) handleDownloadRequest(d *DownloadRequest) {
	core.DownloadBook(c.irc, d.Book)

	c.send <- WaitResponse{
		MessageType: WAIT,
		Status:      "Download request received",
	}
}
