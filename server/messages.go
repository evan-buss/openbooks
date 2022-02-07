package server

import (
	"encoding/json"
	"strconv"

	"github.com/evan-buss/openbooks/core"
)

type MessageType int

// Available commands. These are sent via integers starting at 1
const (
	STATUS MessageType = iota
	CONNECT
	SEARCH
	DOWNLOAD
)

type NotificationType int

const (
	NOTIFY NotificationType = iota
	SUCCESS
	WARNING
	DANGER
)

func messageToString(s MessageType) string {
	name := []string{"STATUS", "CONNECT", "SEARCH", "DOWNLOAD"}
	i := uint8(s)
	switch {
	case i <= uint8(DOWNLOAD):
		return name[i]
	default:
		return strconv.Itoa(int(i))
	}
}

type Message struct {
	MessageType MessageType `json:"type"`
}

type StatusResponse struct {
	Message
	NotificationType NotificationType `json:"status"`
	Title            string           `json:"title"`
	Detail           string           `json:"detail"`
}

// Request in a generic structure for all requests from the websocket client
type Request struct {
	Message
	Payload json.RawMessage `json:"payload"`
}

type Response struct {
	MessageType MessageType `json:"type"`
	Status      string      `json:"status"`
}

// ErrorResponse is a response sent when something goes wrong (ie. bad JSON parse)
type ErrorResponse struct {
	Error   int    `json:"error"`
	Details string `json:"details"`
}

// ConnectionRequest is a request to start the IRC server
type ConnectionRequest struct{}

// ConnectionResponse
type ConnectionResponse struct {
	StatusResponse
	Name string `json:"name"`
}

// SearchRequest is a request that sends a search request to the IRC server for a specific query
type SearchRequest struct {
	Query string `json:"query"`
}

// SearchResponse is a response that is sent containing BookDetails objects that matched the query
type SearchResponse struct {
	Message
	Books  []core.BookDetail `json:"books"`
	Errors []core.ParseError `json:"errors"`
}

// DownloadRequest is a request to download a specific book from the IRC server
type DownloadRequest struct {
	Book string `json:"book"`
}

// DownloadResponse is a response that sends the requested book to the client
type DownloadResponse struct {
	Message
	Name         string `json:"name"`
	DownloadLink string `json:"downloadLink"`
}

func newStatusResponse(notificationType NotificationType) StatusResponse {
	response := StatusResponse{
		Message:          Message{MessageType: STATUS},
		NotificationType: notificationType,
	}
	return response
}
