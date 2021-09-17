package core

import (
	"bufio"
	"context"
	"log"
	"strings"

	"github.com/evan-buss/openbooks/irc"
)

type event int

const (
	noOp           = event(0)
	Message        = event(1)
	SearchResult   = event(2)
	BookResult     = event(3)
	NoResults      = event(4)
	BadServer      = event(5)
	SearchAccepted = event(6)
	MatchesFound   = event(7)
	ServerList     = event(8)
	Ping           = event(9)
)

// Unique identifiers found in the message for various different events.
const (
	pingMessage            = "PING"
	sendMessage            = "DCC SEND"
	noticeMessage          = "NOTICE"
	noResults              = "Sorry"
	serverUnavailable      = "try another server"
	searchAccepted         = "has been accepted"
	searchResultIdentifier = "_results_for"
	numMatches             = "matches"
	beginUserList          = "353"
	endUserList            = "366"
)

type HandlerFunc func(text string)
type EventHandler map[event]HandlerFunc

func StartReader(ctx context.Context, irc *irc.Conn, handler EventHandler) {
	var users strings.Builder
	scanner := bufio.NewScanner(irc)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
			text := scanner.Text()
			if err := scanner.Err(); err != nil {
				log.Println(err)
			}

			// Send raw message if they want to recieve it (logging purposes)
			if invoke, ok := handler[Message]; ok {
				invoke(text)
			}

			event := noOp
			if strings.Contains(text, sendMessage) {
				if strings.Contains(text, searchResultIdentifier) {
					event = SearchResult
				} else {
					event = BookResult
				}
			} else if strings.Contains(text, noticeMessage) {
				if strings.Contains(text, noResults) {
					event = NoResults
				} else if strings.Contains(text, serverUnavailable) {
					event = BadServer
				} else if strings.Contains(text, searchAccepted) {
					event = SearchAccepted
				} else if strings.Contains(text, numMatches) {
					start := strings.LastIndex(text, "returned") + 9
					end := strings.LastIndex(text, "matches") - 1
					text = text[start:end]
					event = MatchesFound
				}
			} else if strings.Contains(text, beginUserList) {
				users.WriteString(text)
			} else if strings.Contains(text, endUserList) {
				event = ServerList
				text = users.String()
				users.Reset()
			} else if strings.Contains(text, pingMessage) {
				event = Ping
			}

			if invoke, ok := handler[event]; ok {
				go invoke(text)
			}
		}
	}
}
