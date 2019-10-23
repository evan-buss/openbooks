package core

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/evan-buss/openbooks/irc"
)

// ReaderHandler handles and responds to different IRC events
// Both the CLI and Server versions implement this interface
type ReaderHandler interface {
	DownloadSearchResults(text string)
	DownloadBookFile(text string)
	NoResults()
	BadServer()
	SearchAccepted()
	MatchesFound(num string)
}

// Possible messages that are sent by the server. We respond accordingly
const (
	pingMessage       = "PING"
	sendMessage       = "DCC SEND"
	noticeMessage     = "NOTICE"
	noResults         = "Sorry"
	serverUnavailable = "try another server"
	searchAccepted    = "has been accepted"
	numMatches        = "matches"
	userList          = "353"
	endUserList       = "366"
)

// Servers contains the cache of available download servers.
var serverCache ServerCache

// ReadDaemon is designed to be launched as a goroutine. Listens for
// specific messages and dispatches appropriate handler functions
// Params: irc - IRC connection
//				 handler - domain specific handler that responds to IRC events
func ReadDaemon(irc *irc.Conn, handler ReaderHandler) {

	var f *os.File
	var err error
	serverCache = ServerCache{Servers: []string{}, Time: time.Now()}

	if irc.Logging {
		f, err = os.OpenFile("irc_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
		f.WriteString("\n==================== NEW LOG ======================\n")

		if err != nil {
			panic(err)
		}
	}

	// Keep a list of users. We want to accumulate all users and then call the handle users method
	var users strings.Builder
	scanner := bufio.NewScanner(irc)

	for scanner.Scan() {
		text := scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Println("Scanner errror: ", err)
		}

		if irc.Logging {
			f.WriteString(text)
		}

		// Respond to Direct Client-to-Client downloads
		if strings.Contains(text, sendMessage) {
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
		} else if strings.Contains(text, userList) {
			users.WriteString(text) // Accumulate the user list
		} else if strings.Contains(text, endUserList) {
			serverCache.ParseServers(users.String())
			users.Reset()
		} else if strings.Contains(text, pingMessage) {
			irc.PONG("irc.irchighway.net")
		}
	}
}
