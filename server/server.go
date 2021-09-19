package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/evan-buss/openbooks/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/rs/cors"
)

type server struct {
	// Shared app configuration
	config *Config

	// Shared data
	repository *Repository

	// Registered clients.
	clients map[uuid.UUID]*Client

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	log *log.Logger
}

// Config contains settings for server
type Config struct {
	Log         bool
	OpenBrowser bool
	Port        string
	UserName    string
	Persist     bool
	DownloadDir string
	Basepath    string
	Server      string
}

func New(config Config) *server {
	return &server{
		repository: NewRepository(),
		config:     &config,
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[uuid.UUID]*Client),
		log:        log.New(os.Stdout, "SERVER: ", log.LstdFlags|log.Lmsgprefix),
	}
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
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "DELETE"},
	}
	router.Use(cors.New(corsConfig).Handler)

	server := New(config)
	routes := server.registerRoutes()

	ctx, cancel := context.WithCancel(context.Background())
	go server.startClientHub(ctx)
	server.registerGracefulShutdown(cancel)
	router.Mount(config.Basepath, routes)

	server.log.Printf("Base Path: %v\n", config.Basepath)
	if config.OpenBrowser {
		browserUrl := "http://127.0.0.1:" + path.Join(config.Port+config.Basepath)
		util.OpenBrowser(browserUrl)
	}

	server.log.Printf("OpenBooks is listening on port %v", config.Port)
	server.log.Fatal(http.ListenAndServe(":"+config.Port, router))
}

// The client hub is to be run in a goroutine and handles management of
// websocket client registrations.
func (s *server) startClientHub(ctx context.Context) {
	for {
		select {
		case client := <-s.register:
			s.clients[client.uuid] = client
		case client := <-s.unregister:
			if _, ok := s.clients[client.uuid]; ok {
				_, cancel := context.WithCancel(client.ctx)
				close(client.send)
				cancel()
				delete(s.clients, client.uuid)
			}
		case <-ctx.Done():
			for _, client := range s.clients {
				_, cancel := context.WithCancel(client.ctx)
				close(client.send)
				cancel()
				delete(s.clients, client.uuid)
			}
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
	err := os.MkdirAll(path.Join(config.DownloadDir, "books"), os.FileMode(0755))
	if err != nil {
		panic(err)
	}
}
