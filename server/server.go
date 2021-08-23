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
)

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

	registerRoutes(hub)

	if config.OpenBrowser {
		browserUrl := "http://127.0.0.1:" + path.Join(config.Port+config.Basepath)
		fmt.Println(browserUrl)
		openbrowser(browserUrl)
	}

	log.Printf("OpenBooks is listening on port %v", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
