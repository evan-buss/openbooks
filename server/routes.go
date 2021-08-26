package server

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/evan-buss/openbooks/irc"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

//go:embed app/dist
var reactClient embed.FS

func (server *server) registerRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Handle("/*", server.staticFilesHandler("app/dist"))
	router.Get("/ws", server.serveWs())
	router.Get("/connections", server.statsHandler())
	return router
}

// serveWs handles websocket requests from the peer.
func (server *server) serveWs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(req *http.Request) bool {
			return true
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		// Each client gets its own connection, so use different usernames by appending count
		name := fmt.Sprintf("%s-%d", server.config.UserName, len(server.hub.clients)+1)
		log.Printf("Connecting to IRC with name %s\n", name)
		client := &Client{
			hub:        server.hub,
			conn:       conn,
			send:       make(chan interface{}, 128),
			disconnect: make(chan struct{}),
			uuid:       uuid.New(),
			irc:        irc.New(name, "OpenBooks - Search and download eBooks"),
		}

		client.hub.register <- client

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go client.writePump()
		go client.readPump()

		log.Printf("%s: Client connected from %s\n", client.uuid, conn.RemoteAddr().String())
	}
}

func (server *server) staticFilesHandler(assetPath string) http.Handler {
	// update the embedded file system's tree so that index.html is at the root
	app, err := fs.Sub(reactClient, assetPath)
	if err != nil {
		log.Fatal(err)
	}

	// strip the predefined base path and serve the static file
	return http.StripPrefix(server.config.Basepath, http.FileServer(http.FS(app)))
}

func (server *server) statsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "There are currently %d active connections.", len(server.hub.clients))
	}
}
