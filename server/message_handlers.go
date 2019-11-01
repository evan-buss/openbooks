package server

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/evan-buss/openbooks/core"
)

// RequestHandler defines a generic handle() method that is called when a specific request type is made
type RequestHandler interface {
	handle()
}

// messageRouter is used to parse the incoming request and respond appropriately
func messageRouter(message Request) {
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
		log.Println(errors.New("unknown message type"))
	}

	err := json.Unmarshal(message.Payload, &obj)
	if err != nil {
		writeJSON(ErrorResponse{
			Error:   message.RequestType,
			Details: err.Error(),
		})
	}
	obj.handle()
}

// handle ConnectionRequests and either connect to the server or do nothing
func (c ConnectionRequest) handle() {
	if !Conn.irc.IsConnected() {
		core.Join(Conn.irc)
		go core.ReadDaemon(Conn.irc, Handler{}) // Start the Read Daemon

		writeJSON(ConnectionResponse{
			MessageType: CONNECT,
			Status:      "You must wait 30 seconds before searching",
			Wait:        30,
		})
		return
	}
	writeJSON(ConnectionResponse{
		MessageType: CONNECT,
		Status:      "IRC Server Ready",
		Wait:        0,
	})
}

// handle SearchRequests and send the query to the book server
func (s SearchRequest) handle() {
	core.SearchBook(Conn.irc, s.Query)

	writeJSON(WaitResponse{
		MessageType: WAIT,
		Status:      "Search request sent",
	})
}

// handle DownloadRequests by sending the request to the book server
func (d DownloadRequest) handle() {
	core.DownloadBook(Conn.irc, d.Book)

	writeJSON(WaitResponse{
		MessageType: WAIT,
		Status:      "Download request received",
	})
}

// handle ServerRequests by sending the currently available book servers
func (s ServersRequest) handle() {
	servers := make(chan []string, 1)
	go core.GetServers(servers)
	results := <-servers

	writeJSON(ServersResponse{
		MessageType: SERVERS,
		Servers:     results,
	})
}
