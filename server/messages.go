package server

import (
	"encoding/json"
	"strconv"

	"github.com/evan-buss/openbooks/core"
)

// Available commands. These are sent via integers starting at 1
const (
	ERROR = iota + 1
	CONNECT
	SEARCH
	DOWNLOAD
	WAIT
	IRCERROR
)

func messageToString(s int) string {
	name := []string{"INVALID", "ERROR", "CONNECT", "SEARCH", "DOWNLOAD", "WAIT", "IRCERROR"}
	i := uint8(s)
	switch {
	case i <= uint8(IRCERROR):
		return name[i]
	default:
		return strconv.Itoa(int(i))
	}
}

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
type ConnectionRequest struct{}

// ConnectionResponse is a response sent upon successful connection to the IRC server
type ConnectionResponse struct {
	MessageType int    `json:"type"`
	Status      string `json:"status"`
	Name        string `json:"name"`
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

// DownloadRequest is a request to download a specific book from the IRC server
type DownloadRequest struct {
	Book string `json:"book"`
}

// DownloadResponse is a response that sends the requested book to the client
type DownloadResponse struct {
	MessageType  int    `json:"type"`
	Name         string `json:"name"`
	DownloadLink string `json:"downloadLink"`
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
