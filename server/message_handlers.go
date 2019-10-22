package server

import (
	"encoding/json"
	"errors"
	"log"
	"time"

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
	go obj.handle()
}

// handle ConnectionRequests and either connect to the server or do nothing
func (c ConnectionRequest) handle() {
	if !ircConn.IsConnected() {
		log.Println("Connecting to IRC")

		core.Join(ircConn)
		go core.ReadDaemon(ircConn, Handler{}) // Start the Read Daemon
		go core.GetUsers(ircConn)              // Get a list of servers

		writeJSON(ConnectionResponse{
			MessageType: CONNECT,
			Status:      "OpenBooks requires a 30 seconds wait period",
			Wait:        30,
		})
		return
	}
	log.Println("You are already connected to the IRC server")
	writeJSON(ConnectionResponse{
		MessageType: CONNECT,
		Status:      "IRC Server Ready",
		Wait:        0,
	})
}

// handle SearchRequests and send the query to the book server
func (s SearchRequest) handle() {
	log.Println("Received SearchRequest for: " + s.Query)

	core.SearchBook(ircConn, s.Query)

	writeJSON(WaitResponse{
		MessageType: WAIT,
		Status:      "Search request sent",
	})
}

// handle DownloadRequests by sending the request to the book server
func (d DownloadRequest) handle() {
	log.Println("Received DownloadRequest: ", d.Book)
	core.DownloadBook(ircConn, d.Book)

	writeJSON(WaitResponse{
		MessageType: WAIT,
		Status:      "Download request received",
	})
}

// handle ServerRequests by sending the currently available book servers
func (s ServersRequest) handle() {
	log.Println("Received ServersRequest")

	oldCache := time.Now().Sub(core.Servers.Time) > time.Minute

	if len(core.Servers.Servers) == 0 || oldCache {
		log.Println("Retrieving serverlist")
		core.GetUsers(ircConn)

		writeJSON(WaitResponse{
			MessageType: WAIT,
			Status:      "Retrieving available book servers",
		})
		return
	}
	log.Println("using cached server list")
	writeJSON(ServersResponse{
		MessageType: SERVERS,
		Servers:     core.Servers.Servers,
	})
}
