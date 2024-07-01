package server

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/evan-buss/openbooks/core"
)

func (s *server) NewIrcEventHandler() core.EventHandler {
	handler := core.EventHandler{}
	handler[core.SearchResult] = s.searchResultHandler(s.config.DownloadDir)
	handler[core.BookResult] = s.bookResultHandler(s.config.DownloadDir, s.config.DisableBrowserDownloads)
	handler[core.NoResults] = s.noResultsHandler
	handler[core.BadServer] = s.badServerHandler
	handler[core.SearchAccepted] = s.searchAcceptedHandler
	handler[core.MatchesFound] = s.matchesFoundHandler
	handler[core.Ping] = s.pingHandler
	handler[core.ServerList] = s.userListHandler(s.repository)
	handler[core.Version] = s.versionHandler(s.config.UserAgent)
	return handler
}

// searchResultHandler downloads from DCC server, parses data, and sends data to client
func (s *server) searchResultHandler(downloadDir string) core.HandlerFunc {
	return func(text string) {
		extractedPath, err := core.DownloadExtractDCCString(filepath.Join(downloadDir, "books"), text, nil)
		if err != nil {
			s.log.Println(err)
			s.send <- newErrorResponse("Error when downloading search results.")
			return
		}

		bookResults, parseErrors, err := core.ParseSearchFile(extractedPath)
		if err != nil {
			s.log.Println(err)
			s.send <- newErrorResponse("Error when parsing search results.")
			return
		}

		if len(bookResults) == 0 && len(parseErrors) == 0 {
			s.noResultsHandler(text)
			return
		}

		// Output all errors so parser can be improved over time
		if len(parseErrors) > 0 {
			s.log.Printf("%d Search Result Parsing Errors\n", len(parseErrors))
			for _, err := range parseErrors {
				s.log.Println(err)
			}
		}

		s.log.Printf("Sending %d search results.\n", len(bookResults))
		s.send <- newSearchResponse(bookResults, parseErrors)

		err = os.Remove(extractedPath)
		if err != nil {
			s.log.Printf("Error deleting search results file: %v", err)
		}
	}
}

// bookResultHandler downloads the book file and sends it over the websocket
func (s *server) bookResultHandler(downloadDir string, disableBrowserDownloads bool) core.HandlerFunc {
	return func(text string) {
		extractedPath, err := core.DownloadExtractDCCString(filepath.Join(downloadDir, "books"), text, nil)
		if err != nil {
			s.log.Println(err)
			s.send <- newErrorResponse("Error when downloading book.")
			return
		}

		s.log.Printf("Sending book entitled '%s'.\n", filepath.Base(extractedPath))
		s.send <- newDownloadResponse(extractedPath, disableBrowserDownloads)
	}
}

// NoResults is called when the server returns that nothing was found for the query
func (s *server) noResultsHandler(_ string) {
	s.send <- newErrorResponse("No results found for the query.")
}

// BadServer is called when the requested download fails because the server is not available
func (s *server) badServerHandler(_ string) {
	s.send <- newErrorResponse("Server is not available. Try another one.")
}

// SearchAccepted is called when the user's query is accepted into the search queue
func (s *server) searchAcceptedHandler(_ string) {
	s.send <- newStatusResponse(NOTIFY, "Search accepted into the queue.")
}

// MatchesFound is called when the server finds matches for the user's query
func (s *server) matchesFoundHandler(num string) {
	s.send <- newStatusResponse(NOTIFY, fmt.Sprintf("Found %s results for your query.", num))
}

func (s *server) pingHandler(serverUrl string) {
	s.irc.Pong(serverUrl)
}

func (s *server) versionHandler(version string) core.HandlerFunc {
	return func(line string) {
		s.log.Printf("Sending CTCP version response: %s", line)
		core.SendVersionInfo(s.irc, line, version)
	}
}

func (s *server) userListHandler(repo *Repository) core.HandlerFunc {
	return func(text string) {
		repo.servers = core.ParseServers(text)
	}
}
