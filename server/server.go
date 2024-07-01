package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/evan-buss/openbooks/core"
	"github.com/evan-buss/openbooks/irc"
	"github.com/evan-buss/openbooks/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/r3labs/sse/v2"
	"github.com/rs/cors"
)

type server struct {
	// Shared app configuration
	config *Config

	// Shared data
	repository *Repository

	sse *sse.Server

	// Send messages to connected clients
	send chan interface{}

	// Keep track of connected clients
	connectedClients atomic.Int32
	status           chan int

	// Single IRC connection shared among all clients
	irc *irc.Conn

	log *log.Logger

	// Mutex to guard the lastSearch timestamp
	lastSearchMutex sync.Mutex

	// The time the last search was performed. Used to rate limit searches.
	lastSearch time.Time
}

// Config contains settings for server
type Config struct {
	Log                     bool
	Port                    string
	UserName                string
	Persist                 bool
	DownloadDir             string
	Basepath                string
	Server                  string
	EnableTLS               bool
	SearchTimeout           time.Duration
	SearchBot               string
	DisableBrowserDownloads bool
	UserAgent               string
}

func New(config Config) *server {
	server := &server{
		config:           &config,
		repository:       NewRepository(),
		send:             make(chan interface{}, 128),
		connectedClients: atomic.Int32{},
		status:           make(chan int),
		irc:              irc.New(config.UserName, config.UserAgent),
		log:              log.New(os.Stdout, "SERVER: ", log.LstdFlags|log.Lmsgprefix),
		lastSearchMutex:  sync.Mutex{},
		lastSearch:       time.Time{},
	}

	sseServer := sse.New()
	sseServer.AutoReplay = false
	sseServer.AutoStream = true
	sseServer.OnSubscribe = func(streamID string, sub *sse.Subscriber) {
		server.status <- 1
	}
	sseServer.OnUnsubscribe = func(streamID string, sub *sse.Subscriber) {
		server.status <- -1
	}

	server.sse = sseServer

	return server
}

// Start instantiates the web server and opens the browser
func Start(config Config) {
	createBooksDirectory(config)
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)

	corsConfig := cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE"},
	}
	router.Use(cors.New(corsConfig).Handler)

	server := New(config)
	routes := server.registerRoutes()

	ctx, cancel := context.WithCancel(context.Background())
	go server.startEventForwarder(ctx)
	server.registerGracefulShutdown(cancel)
	router.Mount(config.Basepath, routes)

	server.log.Printf("Base Path: %s\n", config.Basepath)
	server.log.Printf("OpenBooks is listening on port %v", config.Port)
	server.log.Printf("Download Directory: %s\n", config.DownloadDir)
	server.log.Printf("Open http://localhost:%v%s in your browser.", config.Port, config.Basepath)
	server.log.Fatal(http.ListenAndServe(":"+config.Port, router))
}

func (server *server) startEventForwarder(ctx context.Context) {
	var destructTimer *time.Timer

	for {
		select {
		case message, ok := <-server.send:
			if !ok {
				continue
			}

			byteMessage, err := json.Marshal(message)
			if err != nil {
				server.log.Printf("Error marshalling message to JSON: %s\n", err)
				return
			}
			server.sse.Publish("events", &sse.Event{
				Data: byteMessage,
			})

		case delta := <-server.status:
			server.connectedClients.Add(int32(delta))

			if server.connectedClients.Load() == 0 {
				// Keep the client and IRC connection alive for 3 minutes in case another client connects
				server.log.Println("No clients connected. Waiting 3 minutes before closing connection.")
				destructTimer = time.AfterFunc(3*time.Minute, func() {
					if server.connectedClients.Load() == 0 {
						server.irc.Disconnect()
						server.log.Println("IRC connection closed.")
					}
				})
			} else {
				if destructTimer != nil && destructTimer.Stop() {
					server.log.Println("IRC connection kept alive.")
				}

				server.connectIRC(ctx)
			}

		case <-ctx.Done():
			server.irc.Disconnect()
			return
		}
	}
}

func (server *server) registerGracefulShutdown(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		server.log.Println("Graceful shutdown.")
		// Close the shutdown channel. Triggering all reader/writer WS handlers to close.
		cancel()
		time.Sleep(time.Second)
		os.Exit(0)
	}()
}

func createBooksDirectory(config Config) {
	err := os.MkdirAll(filepath.Join(config.DownloadDir, "books"), os.FileMode(0755))
	if err != nil {
		panic(err)
	}
}

func (server *server) connectIRC(ctx context.Context) {
	// The IRC connection is re-used if already connected
	if server.irc.IsConnected() {
		server.send <- ConnectionResponse{
			StatusResponse: StatusResponse{
				MessageType:      CONNECT,
				NotificationType: NOTIFY,
				Title:            "Welcome back, re-using open IRC connection.",
				Detail:           fmt.Sprintf("IRC username %s", server.irc.Username),
			},
			Name: server.irc.Username,
		}

		return
	}

	server.log.Println("Connecting to IRC server.")

	err := core.Join(server.irc, server.config.Server, server.config.EnableTLS)
	if err != nil {
		server.log.Println(err)
		server.send <- newErrorResponse("Unable to connect to IRC server.")
		return
	}

	handler := server.NewIrcEventHandler()

	if server.config.Log {
		logger, _, err := util.CreateLogFile(server.irc.Username, server.config.DownloadDir)
		if err != nil {
			server.log.Println(err)
		}
		handler[core.Message] = func(text string) { logger.Println(text) }
	}

	go core.StartReader(ctx, server.irc, handler)

	server.send <- ConnectionResponse{
		StatusResponse: StatusResponse{
			MessageType:      CONNECT,
			NotificationType: SUCCESS,
			Title:            "Welcome, connection established.",
			Detail:           fmt.Sprintf("IRC username %s", server.irc.Username),
		},
		Name: server.irc.Username,
	}
}
