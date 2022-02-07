package server

import (
	"os"
	"path"
	"path/filepath"

	"github.com/evan-buss/openbooks/core"
)

func (server *server) NewIrcEventHandler(client *Client) core.EventHandler {
	handler := core.EventHandler{}
	handler[core.SearchResult] = client.searchResultHandler(server.config.DownloadDir)
	handler[core.BookResult] = client.bookResultHandler(server.config.DownloadDir)
	handler[core.NoResults] = client.noResultsHandler
	handler[core.BadServer] = client.badServerHandler
	handler[core.SearchAccepted] = client.searchAcceptedHandler
	handler[core.MatchesFound] = client.matchesFoundHandler
	handler[core.Ping] = client.pingHandler
	handler[core.ServerList] = client.userListHandler(server.repository)
	handler[core.Version] = client.versionHandler
	return handler
}

// searchResultHandler downloads from DCC server, parses data, and sends data to client
func (c *Client) searchResultHandler(downloadDir string) core.HandlerFunc {
	return func(text string) {
		extractedPath, err := core.DownloadExtractDCCString(filepath.Join(downloadDir, "books"), text, nil)
		if err != nil {
			c.log.Println(err)
		}

		results, errors := core.ParseSearchFile(extractedPath)
		// Output all errors so parser can be improved over time
		if len(errors) > 0 {
			c.log.Printf("%d Search Result Parsing Errors\n", len(errors))
			for _, err := range errors {
				c.log.Println(err)
			}
		}

		if len(results) == 0 && len(errors) == 0 {
			c.noResultsHandler(text)
			return
		}

		c.log.Printf("Sending %d search results.\n", len(results))
		c.send <- SearchResponse{
			Message: Message{MessageType: SEARCH},
			Books:   results,
			Errors:  errors,
		}

		err = os.Remove(extractedPath)
		if err != nil {
			c.log.Printf("Error deleting search results file: %v", err)
		}
	}
}

// bookResultHandler downloads the book file and sends it over the websocket
func (c *Client) bookResultHandler(downloadDir string) core.HandlerFunc {
	return func(text string) {
		extractedPath, err := core.DownloadExtractDCCString(filepath.Join(downloadDir, "books"), text, nil)
		if err != nil {
			c.log.Println(err)

			response := newStatusResponse(DANGER)
			response.Title = err.Error()
			c.send <- response
			return
		}

		fileName := filepath.Base(extractedPath)

		c.log.Printf("Sending book entitled '%s'.\n", fileName)
		c.send <- DownloadResponse{
			Message:      Message{DOWNLOAD},
			Name:         fileName,
			DownloadLink: path.Join("library", fileName),
		}
	}
}

// NoResults is called when the server returns that nothing was found for the query
func (c *Client) noResultsHandler(_ string) {
	response := newStatusResponse(DANGER)
	response.Title = "No results found for the query."
	c.send <- response
}

// BadServer is called when the requested download fails because the server is not available
func (c *Client) badServerHandler(_ string) {
	response := newStatusResponse(DANGER)
	response.Title = "Server is not available. Try another one."
	c.send <- response
}

// SearchAccepted is called when the user's query is accepted into the search queue
func (c *Client) searchAcceptedHandler(_ string) {
	response := newStatusResponse(NOTIFY)
	response.Title = "Search accepted into the queue."
	c.send <- response
}

// MatchesFound is called when the server finds matches for the user's query
func (c *Client) matchesFoundHandler(num string) {
	response := newStatusResponse(NOTIFY)
	response.Title = "Found " + num + " results for your query."
	c.send <- response
}

func (c *Client) pingHandler(serverUrl string) {
	c.irc.Pong(serverUrl)
}

func (c *Client) versionHandler(line string) {
	core.SendVersionInfo(c.irc, line)
}

func (c *Client) userListHandler(repo *Repository) core.HandlerFunc {
	return func(text string) {
		repo.servers = core.ParseServers(text)
	}
}
