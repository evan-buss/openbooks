package server

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rakyll/statik/fs"

	// Load the static content
	_ "github.com/evan-buss/openbooks/server/statik"
)

var ircName string
var numConnections *int32 = new(int32)

// Start instantiates the web server and opens the browser
func Start(userName, port string) {

	ircName = userName
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

	staticFs, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(staticFs))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	// openbrowser("http://127.0.0.1:" + port + "/")

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
