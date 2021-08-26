package server

import (
	"encoding/json"
	"log"

	"github.com/evan-buss/openbooks/core"
)

// RequestHandler defines a generic handle() method that is called when a specific request type is made
type RequestHandler interface {
	handle(c *Client)
}

// messageRouter is used to parse the incoming request and respond appropriately
func (c *Client) routeMessage(message Request) {
	var obj RequestHandler

	switch message.RequestType {
	case CONNECT:
		obj = new(ConnectionRequest)
	case SEARCH:
		obj = new(SearchRequest)
	case DOWNLOAD:
		obj = new(DownloadRequest)
	case SERVERS:
		obj = new(ServersRequest)
	default:
		log.Printf("%s: Unknown request type received.\n", c.uuid.String())
	}

	err := json.Unmarshal(message.Payload, &obj)
	if err != nil {
		log.Printf("%s: Invalid request payload.\n", c.uuid.String())
		c.send <- ErrorResponse{
			Error:   message.RequestType,
			Details: err.Error(),
		}
	}
	obj.handle(c)
}

// handle ConnectionRequests and either connect to the server or do nothing
func (ConnectionRequest) handle(c *Client) {
	core.Join(c.irc)
	// Start the Read Daemon
	// TODO: Figure out a way to pass in the config object.
	// go core.ReadDaemon(c.irc, config.Log, Handler{Client: c}, c.disconnect)
	go core.ReadDaemon(c.irc, false, Handler{Client: c}, c.disconnect)

	c.send <- ConnectionResponse{
		MessageType: CONNECT,
		Status:      "IRC Server Ready",
		Wait:        0,
	}
}

// handle SearchRequests and send the query to the book server
func (s SearchRequest) handle(c *Client) {
	core.SearchBook(c.irc, s.Query)

	c.send <- WaitResponse{
		MessageType: WAIT,
		Status:      "Search request sent",
	}
}

// handle DownloadRequests by sending the request to the book server
func (d DownloadRequest) handle(c *Client) {
	core.DownloadBook(c.irc, d.Book)

	c.send <- WaitResponse{
		MessageType: WAIT,
		Status:      "Download request received",
	}
}

// handle ServerRequests by sending the currently available book servers
func (s ServersRequest) handle(c *Client) {
	servers := make(chan []string, 1)
	go core.GetServers(servers)
	results := <-servers

	c.send <- ServersResponse{
		MessageType: SERVERS,
		Servers:     results,
	}
}
