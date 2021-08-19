package server

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//go:embed app/dist
var reactClient embed.FS

// Config contains settings for server
type Config struct {
	Log         bool
	OpenBrowser bool
	Port        string
	UserName    string
	Persist     bool
	DownloadDir string
}

var config Config
var numConnections *int32 = new(int32)

// Start instantiates the web server and opens the browser
func Start(conf Config) {
	config = conf

	hub := newHub()
	go hub.run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		// Close the shutdown channel. Triggering all reader/writer WS handlers to close.
		close(hub.shutdown)
		time.Sleep(time.Second)
		os.Exit(1)
	}()

	http.Handle("/", AddRoutePrefix("/app/dist/", http.FileServer(http.FS(reactClient))))

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	http.HandleFunc("/connections", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "There are currently %d active connections.", *numConnections)
	})

	if config.OpenBrowser {
		openbrowser("http://127.0.0.1:" + config.Port + "/")
	}

	log.Printf("OpenBooks is listening on port %v", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}

// AddRoutePrefix adds a prefix to the request.
func AddRoutePrefix(prefix string, h http.Handler) http.Handler {
	if prefix == "" {
		return h
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := prefix + r.URL.Path
		rp := prefix + r.URL.RawPath
		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = p
		r2.URL.RawPath = rp
		h.ServeHTTP(w, r2)
	})
}
