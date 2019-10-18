package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Handler defines a generic message handler. All handlers should implement
// the handle() method. The method either returns a JSON string or an error
type Handler interface {
	handle() (interface{}, error)
}

// Available commands. These are sent via integers starting at 0
const (
	CONNECT  = iota
	SEARCH   = iota
	DOWNLOAD = iota
)

// Message is the generic wrapper for all received messages.
// It contains a message type and a generic payload of data.
//
// We first read the messageType to determine the payload structure
type Message struct {
	MessageType int             `json:"type"`
	Payload     json.RawMessage `json:"payload"`
}

// ErrorMessage is a generic error returned to the client
// It contains the command type that was invalid as well as a hint string
type ErrorMessage struct {
	Error   int    `json:"error"`
	Details string `json:"details"`
}

// Client -> Server

// ConnectionPayload is called when the user wants to connect to the
// download server. It can optionally define a custom username.
// If the username is not specified, OpenBooks will use the system username
type ConnectionPayload struct {
	Name string `json:"name"`
}

// ConnectionResponse defines the structure for a json response to a
// ConnectionPayload
type ConnectionResponse struct {
	MessageType int    `json:"type"`
	Status      string `json:"status"`
}

// SearchPayload represents a user's search for a book
//
// We can abstract advanced features that aren't present in the IRC server or
// terminal app. For example advanced querying.
// {
// 	"search",
// 	{
// 		"query": "Slaughterhouse Five",
// 		"servers": [
// 			"oatmeal",
// 			"pondering42"
// 		],
// 		"format": [
// 			"epub",
// 			"rar",
// 			"mobi"
// 		]
// 	}
// }
type SearchPayload struct {
	Query   string   `json:"query"`
	Servers []string `json:"servers"`
	Formats []string `json:"formats"`
}

// SearchResponse defines the structure that will respond to a
// SearchPayload
type SearchResponse struct {
	MessageType int      `json:"type"`
	Books       []string `json:"books"`
}

// Server -> Client

// Server connected success
// Available ebook servers
// Server error
// Download accepted
// Send download results
// Send book file.

func messageRouter(message Message) (interface{}, error) {
	var obj Handler

	switch message.MessageType {
	case CONNECT:
		obj = new(ConnectionPayload)
	case SEARCH:
		obj = new(SearchPayload)
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
func (c ConnectionPayload) handle() (interface{}, error) {
	return ConnectionResponse{
		MessageType: CONNECT,
		Status:      "succcessful connection :)",
	}, nil
}

func (s SearchPayload) handle() (interface{}, error) {
	fmt.Println(s)
	return SearchResponse{
		MessageType: SEARCH,
		Books: []string{
			"slaughterhouse five",
			"where the red fern grows",
		},
	}, nil
}
