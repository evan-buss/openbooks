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

// Request is the generic wrapper for all received messages.
// It contains a message type and a generic payload of data.
//
// We first read the messageType to determine the payload structure
type Request struct {
	RequestType int             `json:"type"`
	Payload     json.RawMessage `json:"payload"`
}

// ErrorMessage is a generic error returned to the client
// It contains the command type that was invalid as well as a hint string
type ErrorMessage struct {
	Error   int    `json:"error"`
	Details string `json:"details"`
}

// ConnectionRequest is called when the user wants to connect to the
// download server. It can optionally define a custom username.
// If the username is not specified, OpenBooks will use the system username
type ConnectionRequest struct {
	Name string `json:"name"`
}

// ConnectionResponse defines the structure for a json response to a
// ConnectionPayload
type ConnectionResponse struct {
	MessageType int    `json:"type"`
	Status      string `json:"status"`
	Wait        int    `json:"wait"`
}

// SearchRequest represents a user's search for a book
//
// We can abstract advanced features that aren't present in the IRC server or
// terminal app. For example advanced querying.
type SearchRequest struct {
	Query   string   `json:"query"`
	Servers []string `json:"servers"`
	Formats []string `json:"formats"`
}

// SearchResponse defines the structure that will respond to a
// SearchPayload
type SearchResponse struct {
	MessageType int          `json:"type"`
	Books       []BookDetail `json:"books"`
}

// BookDetail is a single book response
type BookDetail struct {
	Server string `json:"server"`
	Author string `json:"author"`
	Title  string `json:"title"`
	Format string `json:"format"`
	Size   string `json:"size"`
	Full   string `json:"full"`
}

// ServersRequest is a request for available servers
type ServersRequest struct{}

// ServersResponse is a response with a list of available servers
type ServersResponse struct {
	MessageType int      `json:"type"`
	Servers     []string `json:"servers"`
}

// DownloadRequest is a request for a specific book string
type DownloadRequest struct {
	Book string `json:"book"`
}

// DownloadResponse is a response with a book file
type DownloadResponse struct {
	MessageType int    `json:"type"`
	Name        string `json:"name"`
	File        []byte `json:"file"`
}

// WaitResponse is sent when the server has successfully received a
// message but nothing is available yet. Such as searching, we don't
// have an exact time frame for when it will complete.
type WaitResponse struct {
	MessageType int    `json:"type"`
	Status      string `json:"status"`
}

// IrcErrorResponse is sent when some internal IRC response means
// that the request cannot be completed as planned
type IrcErrorResponse struct {
	MessageType int    `json:"type"`
	Status      string `json:"status"`
}

// RequestHandler defines a generic message handler. All handlers should implement
// the handle() method. The method either returns a JSON string or an error
type RequestHandler interface {
	handle() (interface{}, error)
}

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

// Sends a connection response object or a connection error object
func (c ConnectionRequest) handle() (interface{}, error) {
	log.Println("CONNECTION REQUEST. STARTING IRC")
	if !IRC.IsConnected() {
		log.Println("not connected so connecting :)")
		core.Join(IRC)
		go core.ReadDaemon(IRC, Handler{})
	} else {
		log.Println("already connected brother")
	}

	return ConnectionResponse{
		MessageType: CONNECT,
		Status:      "Connected",
		Wait:        30,
	}, nil
}

// Search for the given book.
func (s SearchRequest) handle() (interface{}, error) {
	log.Println("Received SearchRequest for: " + s.Query)
	core.SearchBook(IRC, s.Query)
	// Need to return something, but this is asynchronous...
	return WaitResponse{
		MessageType: WAIT,
		Status:      "Search request sent",
	}, nil
}

// TODO: Implement server parsing functionality
func (s ServersRequest) handle() (interface{}, error) {
	log.Println("Received ServersRequest")
	return ServersResponse{
		MessageType: SERVERS,
	}, nil
}

func (d DownloadRequest) handle() (interface{}, error) {
	log.Println("Received DownloadRequest")
	log.Println(d.Book)
	core.DownloadBook(IRC, d.Book)

	// fmt.Println(d)
	return WaitResponse{
		MessageType: WAIT,
		Status:      "Download request received",
	}, nil
}
