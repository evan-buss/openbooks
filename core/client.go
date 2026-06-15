package core

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/evan-buss/openbooks/dcc"
	"github.com/evan-buss/openbooks/irc"
)

// Handlers contains the callbacks an IrcClient invokes as IRC events are
// processed. Post-processing (DCC download + parsing) is performed before these
// are called, so consumers receive ready-to-use data. Nil callbacks are skipped.
//
// This is the single integration point every OpenBooks interface (CLI, web
// server, Discord bot) implements — the IRC search/download orchestration lives
// here once, and each UI only decides what to do with the results.
type Handlers struct {
	// SearchResults is called after a search results file has been downloaded
	// and parsed. resultsFilePath is the path to the (still on-disk) results
	// file; the consumer owns it and should delete it when done if desired.
	SearchResults func(results []BookDetail, parseErrors []ParseError, resultsFilePath string)

	// BookDownloaded is called with the path to a downloaded book file.
	BookDownloaded func(filePath string)

	// NoResults is called when a search returns nothing.
	NoResults func()

	// BadServer is called when a download targets an unavailable server.
	BadServer func()

	// SearchAccepted is called when a search is accepted into the queue.
	SearchAccepted func()

	// MatchesFound is called with the number of matches the server reported.
	MatchesFound func(count string)

	// ServerList is called whenever the channel user list is (re)parsed.
	ServerList func(servers IrcServers)

	// Message is called with every raw IRC line (intended for logging).
	Message func(line string)

	// Erred is called when downloading or parsing a result/book fails.
	Erred func(err error)

	// Progress, if set, returns the io.Writer used to report download progress
	// for a given file. It may return nil.
	Progress func(filename string, size int64) io.Writer
}

// IrcClient wraps an IRC connection with the shared "search for and download
// books" orchestration used by every OpenBooks interface.
type IrcClient struct {
	Conn *irc.Conn

	saveDir    string // directory where downloaded files are written
	version    string // CTCP VERSION reply
	serverName string // used for PONG replies
	servers    IrcServers
}

// NewIrcClient creates a client whose downloaded files are written to saveDir.
// version is reported to the IRC server in response to CTCP VERSION queries.
func NewIrcClient(username, version, saveDir string) *IrcClient {
	return &IrcClient{
		Conn:    irc.New(username, version),
		saveDir: saveDir,
		version: version,
	}
}

// Username returns the IRC nick of this client.
func (c *IrcClient) Username() string { return c.Conn.Username }

// Servers returns the most recently parsed channel user list.
func (c *IrcClient) Servers() IrcServers { return c.servers }

// Connect joins the IRC server and the #ebooks channel.
func (c *IrcClient) Connect(server string, enableTLS bool) error {
	c.serverName = server
	return Join(c.Conn, server, enableTLS)
}

// Disconnect closes the IRC connection.
func (c *IrcClient) Disconnect() { c.Conn.Disconnect() }

// Search sends a query to the search bot. Results are delivered via Handlers.
func (c *IrcClient) Search(searchBot, query string) { SearchBook(c.Conn, searchBot, query) }

// Download requests a book by its full result line. The file is delivered via
// Handlers.
func (c *IrcClient) Download(book string) { DownloadBook(c.Conn, book) }

// StartReader begins processing IRC events in a goroutine, dispatching to the
// given handlers. It returns immediately.
func (c *IrcClient) StartReader(ctx context.Context, h Handlers) {
	handler := EventHandler{}

	handler[Ping] = func(string) { c.Conn.Pong(c.serverName) }
	handler[Version] = func(line string) { SendVersionInfo(c.Conn, line, c.version) }
	handler[ServerList] = func(text string) {
		c.servers = ParseServers(text)
		if h.ServerList != nil {
			h.ServerList(c.servers)
		}
	}
	if h.Message != nil {
		handler[Message] = h.Message
	}

	handler[SearchResult] = func(text string) {
		path, err := DownloadExtractDCCString(c.saveDir, text, c.progressFor(h, text))
		if err != nil {
			c.cleanupTemp(text)
			invoke(h.Erred, err)
			return
		}
		books, parseErrors, err := ParseSearchFile(path)
		if err != nil {
			os.Remove(path)
			invoke(h.Erred, err)
			return
		}
		if len(books) == 0 && len(parseErrors) == 0 {
			os.Remove(path)
			if h.NoResults != nil {
				h.NoResults()
			}
			return
		}
		if h.SearchResults != nil {
			h.SearchResults(books, parseErrors, path)
		}
	}

	handler[BookResult] = func(text string) {
		path, err := DownloadExtractDCCString(c.saveDir, text, c.progressFor(h, text))
		if err != nil {
			c.cleanupTemp(text)
			invoke(h.Erred, err)
			return
		}
		if h.BookDownloaded != nil {
			h.BookDownloaded(path)
		}
	}

	handler[NoResults] = func(string) {
		if h.NoResults != nil {
			h.NoResults()
		}
	}
	handler[BadServer] = func(string) {
		if h.BadServer != nil {
			h.BadServer()
		}
	}
	handler[SearchAccepted] = func(string) {
		if h.SearchAccepted != nil {
			h.SearchAccepted()
		}
	}
	handler[MatchesFound] = func(count string) {
		if h.MatchesFound != nil {
			h.MatchesFound(count)
		}
	}

	go StartReader(ctx, c.Conn, handler)
}

// progressFor resolves the progress writer (if any) for a given DCC string.
func (c *IrcClient) progressFor(h Handlers, dccStr string) io.Writer {
	if h.Progress == nil {
		return nil
	}
	d, err := dcc.ParseString(dccStr)
	if err != nil {
		return nil
	}
	return h.Progress(d.Filename, d.Size)
}

// cleanupTemp removes the partial ".temp" file left behind by a failed transfer.
func (c *IrcClient) cleanupTemp(dccStr string) {
	d, err := dcc.ParseString(dccStr)
	if err != nil {
		return
	}
	os.Remove(filepath.Join(c.saveDir, d.Filename+".temp"))
}

func invoke(fn func(error), err error) {
	if fn != nil {
		fn(err)
	}
}
