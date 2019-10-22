package core

import (
	"fmt"
	"os"
	"strings"

	"github.com/evan-buss/openbooks/irc"
)

// EventHandler handles and responds to different IRC events
// Both the CLI and Server versions implement this interface
type EventHandler interface {
	DownloadSearchResults(text string)
	DownloadBookFile(text string)
	NoResults()
	BadServer()
	SearchAccepted()
	MatchesFound(num string)
}

// Possible messages that are sent by the server. We respond accordingly
const (
	sendMessage       = "DCC SEND"
	noticeMessage     = "NOTICE"
	noResults         = "Sorry"
	serverUnavailable = "try another server"
	searchAccepted    = "has been accepted"
	numMatches        = "matches"
)

// ReadDaemon is designed to be launched as a goroutine. Listens for
//  specific messages and responds accordingly
// Params: irc - IRC connection
//  			 statusC - boolean channel returns true when the download has finished
// 				 stateC - boolean channel recieves messages from user menu
// 									(true = book download, false = search results download)
func ReadDaemon(irc *irc.Conn, handler EventHandler) {

	var f *os.File
	var err error

	if irc.Logging {
		f, err = os.OpenFile("irc_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
		f.WriteString("\n==================== NEW LOG ======================\n")

		if err != nil {
			panic(err)
		}
	}

	for {
		text := irc.GetMessage()

		if irc.Logging {
			f.WriteString(text)
		}

		if strings.Contains(text, sendMessage) {
			// Respond to Direct Client-to-Client downloads

			if strings.Contains(text, "_results_for") {
				fmt.Println("SEARCH RESULTS")
				go handler.DownloadSearchResults(text)
			} else {
				fmt.Println("BOOK DOWNLOAD")
				go handler.DownloadBookFile(text)
			}
		} else if strings.Contains(text, noticeMessage) {
			if strings.Contains(text, noResults) {
				handler.NoResults()
			} else if strings.Contains(text, serverUnavailable) {
				handler.BadServer()
			} else if strings.Contains(text, searchAccepted) {
				handler.SearchAccepted()
			} else if strings.Contains(text, numMatches) {
				start := strings.LastIndex(text, "returned") + 9
				end := strings.LastIndex(text, "matches") - 1
				handler.MatchesFound(text[start:end])
			}
		}
	}
}
