package core

import (
	"fmt"
	"os"
	"strings"

	"github.com/evan-buss/openbooks/dcc"
	"github.com/evan-buss/openbooks/irc"
)

// ReadDaemon is designed to be launched as a goroutine. Listens for
//  specific messages and responds accordingly
// Params: irc - IRC connection
//  			 statusC - boolean channel returns true when the download has finished
// 				 stateC - boolean channel recieves messages from user menu
// 									(true = book download, false = search results download)
func ReadDaemon(irc *irc.Conn, statusC chan<- bool, stateC <-chan bool) {

	var f *os.File
	var err error

	if irc.Logging {
		f, err = os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
		f.WriteString("\n==================== NEW LOG ======================\n")

		if err != nil {
			panic(err)
		}
	}

	isBook := false

	doneChan := make(chan bool)

	for {
		text := irc.GetMessage()

		if irc.Logging {
			f.WriteString(text)
		}

		select {
		// Get state of app from menu. isBook = download book; !isBook = search
		case isBook = <-stateC:
		default:
		}

		if strings.Contains(text, "DCC SEND") {
			// Respond to Direct Client-to-Client downloads
			go dcc.NewDownload(text, isBook, doneChan)
		} else if strings.Contains(text, "NOTICE") {
			// Notice is a message sent directly to the user
			if strings.Contains(text, "Sorry") {
				// There were no results
				fmt.Println("No results returned for that search...")
				statusC <- !isBook
			} else if strings.Contains(text, "try another server") {
				fmt.Println("That server is not available. Try again...")
				statusC <- !isBook
			} else if strings.Contains(text, "has been accepted") {
				fmt.Println("Search has been accepted. Please wait.")
			} else if strings.Contains(text, "matches") {
				start := strings.LastIndex(text, "returned") + 9
				end := strings.LastIndex(text, "matches") - 1
				fmt.Println("Your search returned " + text[start:end] + " matches.")
			}
		}
		select {
		case <-doneChan:
			// Send message when finished downloading
			statusC <- true
		default:
		}
	}
}
