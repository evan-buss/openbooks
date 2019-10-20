package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// Handler defines a generic message handler. All handlers should implement
// the handle() method. The method either returns a JSON string or an error
type Handler interface {
	handle() (interface{}, error)
}

// ERROR represents some sort of error outside of specific handlers
// such as parsing requests
const ERROR = 0

// Available commands. These are sent via integers starting at 0
const (
	CONNECT = iota + 1
	SEARCH
	DOWNLOAD
	SERVERS
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
	Size   string `json:"size"`
	Format string `json:"format"`
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
	book string `json:"download"`
}

// DownloadResponse is a response with a book file
type DownloadResponse struct {
	MessageType int    `json:"type"`
	Name        string `json:"name"`
	File        []byte `json:"file"`
}

func messageRouter(message Request) (interface{}, error) {
	var obj Handler

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
		return nil, errors.New("Invalid message type")
	}

	err := json.Unmarshal(message.Payload, &obj)
	if err != nil {
		return nil, err
	}
	return obj.handle()
}

// Sends a connection response object or a connection error object
func (c ConnectionRequest) handle() (interface{}, error) {
	return ConnectionResponse{
		MessageType: CONNECT,
		Status:      "Connected",
	}, nil
}

func (s SearchRequest) handle() (interface{}, error) {
	return SearchResponse{
		MessageType: SEARCH,
		Books: []BookDetail{
			BookDetail{
				Server: "Oatmeal",
				Author: "Kurt Vonnegut",
				Title:  "Slaughterhouse Five",
				Format: "PDF",
				Size:   "220.0KB",
				Full:   "THIS IS THE FULL DOWNLOAD STRING",
			},
			BookDetail{
				Server: "Pondering42",
				Author: "Kurt Vonnegut",
				Title:  "Slaughterhouse Five",
				Format: "rar",
				Size:   "245.9KB",
				Full:   "THIS IS THE FULL DOWNLOAD STRING",
			},
		},
	}, nil
}

func (s ServersRequest) handle() (interface{}, error) {
	return ServersResponse{
		MessageType: SERVERS,
	}, nil
}

func (d DownloadRequest) handle() (interface{}, error) {
	fileName := "blood.mobi"
	// fileName := "cover.jpg"
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	// fmt.Println(d)
	return DownloadResponse{
		MessageType: DOWNLOAD,
		Name:        fileName,
		File:        dat,
	}, nil
}
