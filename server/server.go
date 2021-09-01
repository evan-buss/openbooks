package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/evan-buss/openbooks/core"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type server struct {
	router     *chi.Mux
	hub        *Hub
	config     *Config
	repository *Repository
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
}

// Start instantiates the web server and opens the browser
func Start(config Config) {
	config.Basepath = sanitizePath(config.Basepath)
	fmt.Printf("Base Path: %v\n", config.Basepath)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	hub := newHub()
	go hub.run()

	server := &server{
		config:     &config,
		router:     router,
		hub:        hub,
		repository: &Repository{servers: new(core.IrcServers)},
	}

	routes := server.registerRoutes()
	router.Mount(config.Basepath, routes)
	server.registerGracefulShutdown()

	if config.OpenBrowser {
		browserUrl := "http://127.0.0.1:" + path.Join(config.Port+config.Basepath)
		openbrowser(browserUrl)
	}

	log.Printf("OpenBooks is listening on port %v", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, router))
}

func (server *server) registerGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Graceful shutdown.")
		// Close the shutdown channel. Triggering all reader/writer WS handlers to close.
		close(server.hub.shutdown)
		time.Sleep(time.Second)
		os.Exit(1)
	}()
}

func sanitizePath(basepath string) string {
	cleaned := path.Clean(basepath)
	if cleaned == "/" {
		return cleaned
	}
	return cleaned + "/"
}
