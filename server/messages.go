package server

import (
	"encoding/json"
	"fmt"
	"math"
	"path"

	"github.com/evan-buss/openbooks/core"
)

//go:generate stringer -type=MessageType
type MessageType int

// Available commands. These are sent via integers starting at 1
const (
	STATUS MessageType = iota
	CONNECT
	SEARCH
	DOWNLOAD
	RATELIMIT
)

type NotificationType int

const (
	NOTIFY NotificationType = iota
	SUCCESS
	WARNING
	DANGER
)

type StatusResponse struct {
	MessageType      MessageType      `json:"type"`
	NotificationType NotificationType `json:"appearance"`
	Title            string           `json:"title"`
	Detail           string           `json:"detail"`
}

// Request in a generic structure for all requests from the websocket client
type Request struct {
	MessageType MessageType     `json:"type"`
	Payload     json.RawMessage `json:"payload"`
}

// ConnectionRequest is a request to start the IRC server
type ConnectionRequest struct{}

// SearchRequest is a request that sends a search request to the IRC server for a specific query
type SearchRequest struct {
	Query string `json:"query"`
}

// DownloadRequest is a request to download a specific book from the IRC server
type DownloadRequest struct {
	Book   string `json:"book"`
	Author string `json:"author,omitempty"`
	Title  string `json:"title,omitempty"`
}

// ConnectionResponse
type ConnectionResponse struct {
	StatusResponse
	Name string `json:"name"`
}

// SearchResponse is a response that is sent containing BookDetails objects that matched the query
type SearchResponse struct {
	StatusResponse
	Books  []core.BookDetail `json:"books"`
	Errors []core.ParseError `json:"errors"`
}

// DownloadResponse is a response that sends the requested book to the client
type DownloadResponse struct {
	StatusResponse
	Name         string `json:"name"`
	DownloadPath string `json:"downloadPath"`
}

func newRateLimitResponse(remainingSeconds float64) StatusResponse {
	wait := math.Round(remainingSeconds)
	units := "seconds"
	if wait == 1 {
		units = "second"
	}

	return StatusResponse{
		MessageType:      RATELIMIT,
		NotificationType: WARNING,
		Title:            "You are searching too frequently!",
		Detail:           fmt.Sprintf("Please wait %v %s to submit another search.", wait, units),
	}
}

func newSearchResponse(results []core.BookDetail, errors []core.ParseError) SearchResponse {
	detail := fmt.Sprintf("There were %v parsing errors.", len(errors))
	if len(errors) == 1 {
		detail = "There was 1 parsing error."
	}
	return SearchResponse{
		StatusResponse: StatusResponse{
			MessageType:      SEARCH,
			NotificationType: SUCCESS,
			Title:            fmt.Sprintf("%v Search Results Received", len(results)),
			Detail:           detail,
		},
		Books:  results,
		Errors: errors,
	}
}

func newDownloadResponse(filePath string, disableBrowserDownloads bool) DownloadResponse {
	// If we don't want to autodownload the file, show the user the path to the file
	// otherwise just show file name.
	if !disableBrowserDownloads {
		filePath = path.Base(filePath)
	}

	response := DownloadResponse{
		StatusResponse: StatusResponse{
			MessageType:      DOWNLOAD,
			NotificationType: SUCCESS,
			Title:            "Book file received.",
			Detail:           filePath,
		},
	}

	// If we want to autodownload the file, add the path to the response
	// client will not attempt autodownload if the path is empty
	if !disableBrowserDownloads {
		response.DownloadPath = path.Join("library", filePath)
	}

	return response
}

func newStatusResponse(notificationType NotificationType, title string) StatusResponse {
	return StatusResponse{
		MessageType:      STATUS,
		NotificationType: notificationType,
		Title:            title,
	}
}

func newErrorResponse(title string) StatusResponse {
	return StatusResponse{
		MessageType:      STATUS,
		NotificationType: DANGER,
		Title:            title,
	}
}
