package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

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
	SearchTimeout           time.Duration
	SearchBot               string
	DisableBrowserDownloads bool
	UserAgent               string
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
		AllowedOrigins:   []string{"http://127.0.0.1:5173"},
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

	server.log.Printf("Base Path: %s\n", config.Basepath)
	server.log.Printf("OpenBooks is listening on port %v", config.Port)
	server.log.Printf("Download Directory: %s\n", config.DownloadDir)
	server.log.Printf("Open http://localhost:%v%s in your browser.", config.Port, config.Basepath)
	server.log.Fatal(http.ListenAndServe(":"+config.Port, router))
}

// The client hub is to be run in a goroutine and handles management of
// websocket client registrations.
func (server *server) startClientHub(ctx context.Context) {
	type selfDestructor struct {
		timer  *time.Timer
		client *Client
	}

	selfDestructors := make(map[uuid.UUID]selfDestructor)

	for {
		select {
		case client := <-server.register:
			if destructor, ok := selfDestructors[client.uuid]; ok {
				destructor.timer.Stop()
				server.log.Printf("Client %s reconnected\n", client.uuid.String())

				// Update the existing client's websocket connection to the new one
				destructor.client.conn = client.conn
				client = destructor.client

				delete(selfDestructors, client.uuid)
			}

			server.clients[client.uuid] = client
			go server.writePump(client)
			go server.readPump(client)

		case client := <-server.unregister:
			// Keep the client and IRC connection alive for 3 minutes in case the client reconnects
			timer := time.AfterFunc(time.Minute*3, func() {
				if _, ok := selfDestructors[client.uuid]; ok {
					client.irc.Disconnect()
					_, cancel := context.WithCancel(client.ctx)
					close(client.send)
					cancel()
					delete(selfDestructors, client.uuid)
					server.log.Printf("Client %s self-destructed\n", client.uuid.String())
				}
			})

			selfDestructors[client.uuid] = selfDestructor{timer, client}
			delete(server.clients, client.uuid)
			server.log.Printf("Client %s disconnected\n", client.uuid.String())
		case <-ctx.Done():
			for _, client := range server.clients {
				_, cancel := context.WithCancel(client.ctx)
				close(client.send)
				cancel()
				delete(server.clients, client.uuid)
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
	err := os.MkdirAll(filepath.Join(config.DownloadDir, "books"), os.FileMode(0755))
	if err != nil {
		panic(err)
	}
}
