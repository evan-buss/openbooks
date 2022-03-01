package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/robfig/cron/v3"
)

type Config struct {
	Username        string            `mapstructure:"username"`
	Port            string            `mapstructure:"port"`
	IrcServers      []string          `mapstructure:"ircServers"`
	RefreshSchedule string            `mapstructure:"refreshSchedule"`
	AdminUsers      map[string]string `mapstructure:"adminusers"`
}

type server struct {
	searcher Searcher
	indexer  Indexer
	router   *chi.Mux
	cron     *cron.Cron
	config   Config
}

func NewServer(config Config) {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	// router.Use(cors.New(cors.Options{
	// 	AllowCredentials: true,
	// 	AllowedOrigins:   []string{"http://localhost:3000"},
	// 	AllowedHeaders:   []string{"*"},
	// 	AllowedMethods:   []string{"GET", "DELETE"},
	// }).Handler)

	searcher, err := NewBleveSeacher()
	if err != nil {
		log.Fatal(err)
	}

	server := &server{
		searcher: searcher,
		indexer:  &DirectoryIndexer{"sources"},
		router:   router,
		cron:     cron.New(),
		config:   config,
	}

	server.registerRoutes()
	server.cron.AddFunc(server.config.RefreshSchedule, server.RefreshJob)
	server.cron.Start()

	httpServer := &http.Server{
		Addr:         ":" + config.Port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      server.router,
	}

	log.Println("Listening on ", httpServer.Addr)
	log.Fatal(httpServer.ListenAndServe())
}
