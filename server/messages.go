package server

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/evan-buss/openbooks/core"
)

// ERROR represents some sort of error outside of specific handlers
// such as parsing requests
const ERROR = 0

// Available commands. These are sent via integers starting at 0
const (
	CONNECT = iota + 1
	SEARCH
	DOWNLOAD
	SERVERS
	WAIT
	IRCERROR
)

// Request in a generic structure for all requests from the websocket client
type Request struct {
	RequestType int             `json:"type"`
	Payload     json.RawMessage `json:"payload"`
}

// ErrorResponse is a response sent when something goes wrong (ie. bad JSON parse)
type ErrorResponse struct {
	Error   int    `json:"error"`
	Details string `json:"details"`
}

// ConnectionRequest is a request to start the IRC server
type ConnectionRequest struct {
	Name string `json:"name"`
}

// TODO: Make the client read from the Wait variable and set the timer accordingly
// ConnectionResponse is a response sent upon successful connection to the IRC server
type ConnectionResponse struct {
	MessageType int    `json:"type"`
	Status      string `json:"status"`
	Wait        int    `json:"wait"`
}

// SearchRequest is a request that sends a search request to the IRC server for a specific query
type SearchRequest struct {
	Query string `json:"query"`
}

// SearchResponse is a response that is sent containing BookDetails objects that matched the query
type SearchResponse struct {
	MessageType int          `json:"type"`
	Books       []BookDetail `json:"books"`
}

// BookDetail contains the details of a single Book found on the IRC server
type BookDetail struct {
	Server string `json:"server"`
	Author string `json:"author"`
	Title  string `json:"title"`
	Format string `json:"format"`
	Size   string `json:"size"`
	Full   string `json:"full"`
}

// ServersRequest is a request that lists available IRC servers
type ServersRequest struct{}

// ServersResponse is a response that lists the IRC servers that are online and available
type ServersResponse struct {
	MessageType int      `json:"type"`
	Servers     []string `json:"servers"`
}

// DownloadRequest is a request to download a specific book from the IRC server
type DownloadRequest struct {
	Book string `json:"book"`
}

// DownloadResponse is a response that sends the requested book to the client
type DownloadResponse struct {
	MessageType int    `json:"type"`
	Name        string `json:"name"`
	File        []byte `json:"file"`
}

// WaitResponse is a response that reports status updates to the client. IRC is asynchronous
// and has an unbounded time-frame so we want to show the client things are happening
type WaitResponse struct {
	MessageType int    `json:"type"`
	Status      string `json:"status"`
}

// IrcErrorResponse is a response that indicates something went wrong on the IRC server's end
type IrcErrorResponse struct {
	MessageType int    `json:"type"`
	Status      string `json:"status"`
}

// RequestHandler defines a generic handle() method that is called when a specific request type is made
type RequestHandler interface {
	handle() (interface{}, error)
}

// messageRouter is used to parse the incoming request and respond appropriately
func messageRouter(message Request) (interface{}, error) {
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
		return nil, errors.New("invalid message type")
	}

	err := json.Unmarshal(message.Payload, &obj)
	if err != nil {
		return nil, err
	}
	return obj.handle()
}

// handle ConnectionRequests and either connect to the server or do nothing
func (c ConnectionRequest) handle() (interface{}, error) {
	if !IRC.IsConnected() {
		log.Println("Connecting to IRC")
		core.Join(IRC)
		go core.ReadDaemon(IRC, Handler{})
		return ConnectionResponse{
			MessageType: CONNECT,
			Status:      "IRC Server Requires 30 second wait period",
			Wait:        30,
		}, nil
	} else {
		log.Println("You are already connected to the IRC server")
		return ConnectionResponse{
			MessageType: CONNECT,
			Status:      "IRC Server Ready",
			Wait:        0,
		}, nil
	}
}

// handle SearchRequests and send the query to the IRC server
func (s SearchRequest) handle() (interface{}, error) {
	log.Println("Received SearchRequest for: " + s.Query)

	core.SearchBook(IRC, s.Query)

	// Update client that request was sent
	return WaitResponse{
		MessageType: WAIT,
		Status:      "Search request sent",
	}, nil
}

// handle ServerRequests by sending the currently available book servers
func (s ServersRequest) handle() (interface{}, error) {
	log.Println("Received ServersRequest")
	core.GetUsers(IRC)

	//TODO: Implement server listing and parsing.
	return WaitResponse{
		MessageType: SERVERS,
		Status:      "Retrieving available book servers",
	}, nil
}

// handle DownloadRequests by sending the request to the book server
func (d DownloadRequest) handle() (interface{}, error) {
	log.Println("Received DownloadRequest: ", d.Book)
	core.DownloadBook(IRC, d.Book)

	return WaitResponse{
		MessageType: WAIT,
		Status:      "Download request received",
	}, nil
}
