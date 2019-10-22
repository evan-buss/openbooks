package server

import (
	"encoding/json"

	"github.com/evan-buss/openbooks/core"
)

// Messages defines the client and server messages that are passed between
// the websocket / website server and the client (website)

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
	MessageType int               `json:"type"`
	Books       []core.BookDetail `json:"books"`
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
