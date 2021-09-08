package core

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

// TODO: Should this take a context.Context for implicit cancellations?
type ReadDaemon struct {
	Reader     io.Reader
	Disconnect <-chan struct{}
	LogConfig  LogConfig
	Events     IrcHighwayEvents
}

type LogConfig struct {
	Enable   bool
	UserName string
	Path     string
}

// IrcHighwayEvents handles and responds to different IRC events
// Both the CLI and Server versions implement this interface
type IrcHighwayEvents interface {
	DownloadSearchResults(text string)
	DownloadBookFile(text string)
	NoResults()
	BadServer()
	SearchAccepted()
	ServerList(servers IrcServers)
	MatchesFound(num string)
	Ping()
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
	beginUserList     = "353"
	endUserList       = "366"
)

// Start is designed to be launched as a goroutine. Listens for
// specific messages and dispatches appropriate handler functions
func (r *ReadDaemon) Start() {
	ircLogger, closer, err := createLogFile(r.LogConfig)
	if err != nil {
		log.Println(err)
		return
	}

	if closer != nil {
		defer closer.Close()
	}

	// User list can be long enough to necessitate multiple messages. IRC defines
	// message codes denoting the beginning and end of the user list.
	// We append to the string.Builder and then flush when we receive the end code.
	var users strings.Builder
	scanner := bufio.NewScanner(r.Reader)
	for scanner.Scan() {
		select {
		case <-r.Disconnect:
			return

		default:
			text := scanner.Text()
			if err := scanner.Err(); err != nil {
				log.Println(err)
			}

			ircLogger.Println(text)

			// Respond to Direct Client-to-Client downloads
			if strings.Contains(text, sendMessage) {
				if strings.Contains(text, "_results_for") {
					go r.Events.DownloadSearchResults(text)
				} else {
					go r.Events.DownloadBookFile(text)
				}
			} else if strings.Contains(text, noticeMessage) {
				if strings.Contains(text, noResults) {
					r.Events.NoResults()
				} else if strings.Contains(text, serverUnavailable) {
					r.Events.BadServer()
				} else if strings.Contains(text, searchAccepted) {
					r.Events.SearchAccepted()
				} else if strings.Contains(text, numMatches) {
					start := strings.LastIndex(text, "returned") + 9
					end := strings.LastIndex(text, "matches") - 1
					r.Events.MatchesFound(text[start:end])
				}
			} else if strings.Contains(text, beginUserList) {
				users.WriteString(text)
			} else if strings.Contains(text, endUserList) {
				r.Events.ServerList(ParseServers(users.String()))
				users.Reset()
			} else if strings.Contains(text, pingMessage) {
				r.Events.Ping()
			}
		}
	}
}

// createLogFile returns a loger for persisting IRC history. If enable is true,
// the messages are logged to a file in the "logs" folder of the provided path.
// If logging is disabled, the logs aren't written anywhere, but the returned
// logger can be called like normal. Check the returned io.Closer for nil and
// call Close() to clean up resources after you are done writting to the logger.
func createLogFile(config LogConfig) (*log.Logger, io.Closer, error) {
	date := time.Now().Format("2006-01-02--15-04-05")
	fileName := fmt.Sprintf("%s--%s.log", config.UserName, date)

	err := os.MkdirAll(path.Join(config.Path, "logs"), os.FileMode(0755))
	if err != nil {
		return nil, nil, err
	}

	path := path.Join(config.Path, "logs", fileName)
	logFile, err := os.Create(path)
	if err != nil {
		return nil, nil, err
	}

	if config.Enable {
		return log.New(logFile, "", 0), logFile, nil
	}

	return log.New(io.Discard, "", 0), nil, nil
}
