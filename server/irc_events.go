package server

import (
	"os"
	"path"

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
	return handler
}

// searchResultHandler downloads from DCC server, parses data, and sends data to client
func (c *Client) searchResultHandler(downloadDir string) core.HandlerFunc {
	return func(text string) {
		extractedPath, err := core.DownloadExtractDCCString(path.Join(downloadDir, "books"), text, nil)
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

		if len(results) == 0 {
			c.noResultsHandler(text)
			return
		}

		c.log.Printf("Sending %d search results.\n", len(results))
		c.send <- SearchResponse{
			MessageType: SEARCH,
			Books:       results,
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
		extractedPath, err := core.DownloadExtractDCCString(path.Join(downloadDir, "books"), text, nil)
		if err != nil {
			c.log.Println(err)
			c.send <- ErrorResponse{
				Error:   ERROR,
				Details: err.Error(),
			}
			return
		}

		fileName := path.Base(extractedPath)

		c.log.Printf("Sending book entitled '%s'.\n", fileName)
		c.send <- DownloadResponse{
			MessageType:  DOWNLOAD,
			Name:         fileName,
			DownloadLink: path.Join("/static/library", fileName),
		}
	}
}

// NoResults is called when the server returns that nothing was found for the query
func (c *Client) noResultsHandler(_ string) {
	c.send <- IrcErrorResponse{
		MessageType: IRCERROR,
		Status:      "No results found for the query.",
	}
}

// BadServer is called when the requested download fails because the server is not available
func (c *Client) badServerHandler(_ string) {
	c.send <- IrcErrorResponse{
		MessageType: IRCERROR,
		Status:      "Server is not available. Try another one.",
	}
}

// SearchAccepted is called when the user's query is accepted into the search queue
func (c *Client) searchAcceptedHandler(_ string) {
	c.send <- WaitResponse{
		MessageType: WAIT,
		Status:      "Search accepted into the queue.",
	}
}

// MatchesFound is called when the server finds matches for the user's query
func (c *Client) matchesFoundHandler(num string) {
	c.send <- WaitResponse{
		MessageType: WAIT,
		Status:      "Found " + num + " results for your query.",
	}
}

func (c *Client) pingHandler(serverUrl string) {
	c.irc.Pong(serverUrl)
}

func (c *Client) userListHandler(repo *Repository) core.HandlerFunc {
	return func(text string) {
		repo.servers = core.ParseServers(text)
	}
}
