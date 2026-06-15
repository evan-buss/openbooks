package discord

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/evan-buss/openbooks/core"
	"github.com/evan-buss/openbooks/util"
)

// Sentinel errors returned by the broker. They are safe to surface directly to
// users in Discord.
var (
	ErrNoResults = errors.New("no results found for your query")
	ErrBadServer = errors.New("that server is not available, try a different result")
	ErrTimeout   = errors.New("timed out waiting for a response from IRC")
)

const (
	// IRC searches are queued server-side and can take a while during busy
	// periods. Downloads stream a file and are usually quick, but large books
	// or slow servers warrant a generous window.
	searchTimeout   = 5 * time.Minute
	downloadTimeout = 10 * time.Minute
)

// Broker turns the asynchronous, single shared IRC connection into a
// serialized, synchronous request/response API suitable for a multi-user
// Discord bot. The actual IRC/DCC/search work is delegated to core.IrcClient —
// the broker only adds serialization, rate limiting, and result correlation.
type Broker struct {
	client    *core.IrcClient
	searchBot string
	rateLimit time.Duration
	debug     bool

	handlers  core.Handlers
	logCloser io.Closer

	// opMu serializes whole operations end-to-end so that asynchronous IRC
	// events can be unambiguously matched to the waiting caller.
	opMu sync.Mutex

	// sinkMu guards the currently active result sinks.
	sinkMu     sync.Mutex
	searchSink chan []core.BookDetail
	bookSink   chan string
	errSink    chan error

	lastSearch time.Time
}

// NewBroker connects to IRC, joins #ebooks, prepares the download directory, and
// (optionally) opens a raw IRC log file. Call Start to begin processing events.
func NewBroker(cfg Config) (*Broker, error) {
	saveDir := filepath.Join(cfg.DownloadDir, "books")
	if err := os.MkdirAll(saveDir, 0o755); err != nil {
		return nil, fmt.Errorf("creating download directory: %w", err)
	}

	client := core.NewIrcClient(cfg.UserName, cfg.Version, saveDir)
	if err := client.Connect(cfg.Server, cfg.EnableTLS); err != nil {
		return nil, err
	}

	b := &Broker{
		client:    client,
		searchBot: cfg.SearchBot,
		rateLimit: cfg.RateLimit,
		debug:     cfg.Debug,
	}

	handlers := core.Handlers{
		SearchResults: func(books []core.BookDetail, _ []core.ParseError, path string) {
			os.Remove(path)
			b.deliverSearch(books)
		},
		BookDownloaded: func(path string) {
			log.Printf("Transfer complete: %s", path)
			b.deliverBook(path)
		},
		NoResults: func() { b.deliverErr(ErrNoResults) },
		BadServer: func() { b.deliverErr(ErrBadServer) },
		Erred:     func(err error) { b.deliverErr(err) },
		Progress: func(filename string, size int64) io.Writer {
			log.Printf("Receiving %q (%s)", filename, humanBytes(size))
			if b.debug {
				return newProgressLogger(filename, size)
			}
			return nil
		},
	}

	// When logging is enabled, write every raw IRC line to a file — the same
	// firehose (and destination) the CLI/server use. Keeps channel chatter out
	// of stdout.
	if cfg.Log {
		logger, closer, err := util.CreateLogFile(cfg.UserName, cfg.DownloadDir)
		if err != nil {
			log.Printf("Could not create IRC log file: %v", err)
		} else {
			b.logCloser = closer
			handlers.Message = func(text string) { logger.Println(text) }
			log.Printf("Logging raw IRC traffic to %s", filepath.Join(cfg.DownloadDir, "logs"))
		}
	}

	b.handlers = handlers
	return b, nil
}

// Start launches the IRC reader loop. It returns immediately.
func (b *Broker) Start(ctx context.Context) {
	b.client.StartReader(ctx, b.handlers)
}

// Disconnect closes the IRC connection and any open log file.
func (b *Broker) Disconnect() {
	b.client.Disconnect()
	if b.logCloser != nil {
		b.logCloser.Close()
	}
}

// Search sends a query and returns the parsed results. Calls are serialized and
// rate-limited.
func (b *Broker) Search(ctx context.Context, query string) ([]core.BookDetail, error) {
	b.opMu.Lock()
	defer b.opMu.Unlock()

	// Respect the minimum interval between searches.
	if wait := time.Until(b.lastSearch.Add(b.rateLimit)); wait > 0 {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(wait):
		}
	}

	sink := make(chan []core.BookDetail, 1)
	errSink := make(chan error, 1)
	b.setSinks(sink, nil, errSink)
	defer b.clearSinks()

	b.client.Search(b.searchBot, query)
	b.lastSearch = time.Now()
	if b.debug {
		log.Printf("[debug] sent search for %q to @%s", query, b.searchBot)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errSink:
		return nil, err
	case <-time.After(searchTimeout):
		return nil, ErrTimeout
	case books := <-sink:
		return books, nil
	}
}

// Download requests a book (by its full result line) and returns the path to
// the saved file. Calls are serialized.
func (b *Broker) Download(ctx context.Context, bookLine string) (string, error) {
	b.opMu.Lock()
	defer b.opMu.Unlock()

	sink := make(chan string, 1)
	errSink := make(chan error, 1)
	b.setSinks(nil, sink, errSink)
	defer b.clearSinks()

	b.client.Download(bookLine)
	if b.debug {
		log.Printf("[debug] sent download request: %s", bookLine)
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case err := <-errSink:
		return "", err
	case <-time.After(downloadTimeout):
		return "", ErrTimeout
	case path := <-sink:
		return path, nil
	}
}

func (b *Broker) setSinks(search chan []core.BookDetail, book chan string, errs chan error) {
	b.sinkMu.Lock()
	defer b.sinkMu.Unlock()
	b.searchSink, b.bookSink, b.errSink = search, book, errs
}

func (b *Broker) clearSinks() {
	b.setSinks(nil, nil, nil)
}

func (b *Broker) deliverSearch(books []core.BookDetail) {
	b.sinkMu.Lock()
	ch := b.searchSink
	b.sinkMu.Unlock()
	if ch != nil {
		select {
		case ch <- books:
		default:
		}
	}
}

func (b *Broker) deliverBook(path string) {
	b.sinkMu.Lock()
	ch := b.bookSink
	b.sinkMu.Unlock()
	if ch != nil {
		select {
		case ch <- path:
		default:
		}
	}
}

func (b *Broker) deliverErr(err error) {
	b.sinkMu.Lock()
	ch := b.errSink
	b.sinkMu.Unlock()
	if ch != nil {
		select {
		case ch <- err:
		default:
		}
	}
}

// progressLogger is an io.Writer that periodically logs transfer progress.
type progressLogger struct {
	name    string
	total   int64
	written int64
	lastLog time.Time
}

func newProgressLogger(name string, total int64) *progressLogger {
	return &progressLogger{name: name, total: total, lastLog: time.Now()}
}

func (p *progressLogger) Write(b []byte) (int, error) {
	p.written += int64(len(b))
	if time.Since(p.lastLog) >= 2*time.Second || p.written == p.total {
		p.lastLog = time.Now()
		pct := 0.0
		if p.total > 0 {
			pct = 100 * float64(p.written) / float64(p.total)
		}
		log.Printf("[debug] %s: %s / %s (%.0f%%)", p.name, humanBytes(p.written), humanBytes(p.total), pct)
	}
	return len(b), nil
}

func humanBytes(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for x := n / unit; x >= unit; x /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(n)/float64(div), "KMGTPE"[exp])
}
