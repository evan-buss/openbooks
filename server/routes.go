package server

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

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
	router.Get("/stats", server.statsHandler())
	router.Get("/servers", server.serverListHandler())

	router.Group(func(r chi.Router) {
		r.Use(server.requireUser)
		r.Get("/library", server.libraryHandler())
		r.Delete("/library/{fileName}", server.deleteLibraryFileHandler())
		r.Get("/static/library/*", server.libraryFileHandler("/static/library").ServeHTTP)
		r.Get("/test", server.testHandler())
	})

	return router
}

// serveWs handles websocket requests from the peer.
func (server *server) serveWs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("OpenBooks")
		if errors.Is(err, http.ErrNoCookie) {
			cookie = &http.Cookie{
				Name:     "OpenBooks",
				Value:    uuid.New().String(),
				Secure:   true,
				HttpOnly: true,
				Expires:  time.Now().Add(time.Hour * 24 * 7),
				SameSite: http.SameSiteStrictMode,
			}
			w.Header().Add("Set-Cookie", cookie.String())
		}

		userId := uuid.MustParse(cookie.Value)

		upgrader.CheckOrigin = func(req *http.Request) bool {
			return true
		}

		conn, err := upgrader.Upgrade(w, r, w.Header())
		if err != nil {
			server.log.Println(err)
			return
		}

		// Each client gets its own connection, so use different usernames by appending count
		name := fmt.Sprintf("%s-%d", server.config.UserName, len(server.clients)+1)
		client := &Client{
			conn:       conn,
			send:       make(chan interface{}, 128),
			disconnect: make(chan struct{}),
			uuid:       userId,
			irc:        irc.New(name, "OpenBooks - Search and download eBooks"),
			log:        log.New(os.Stdout, fmt.Sprintf("CLIENT (%s): ", name), log.LstdFlags|log.Lmsgprefix),
		}

		server.log.Printf("Client connected from %s\n", conn.RemoteAddr().String())
		client.log.Println("New client created.")

		server.register <- client

		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go server.writePump(client)
		go server.readPump(client)
	}
}

func (server *server) staticFilesHandler(assetPath string) http.Handler {
	// update the embedded file system's tree so that index.html is at the root
	app, err := fs.Sub(reactClient, assetPath)
	if err != nil {
		server.log.Fatal(err)
	}

	// strip the predefined base path and serve the static file
	return http.StripPrefix(server.config.Basepath, http.FileServer(http.FS(app)))
}

func (server *server) statsHandler() http.HandlerFunc {
	type statsReponse struct {
		UUID string `json:"uuid"`
		IP   string `json:"ip"`
		Name string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		result := make([]statsReponse, 0, len(server.clients))

		for _, client := range server.clients {
			details := statsReponse{
				UUID: client.uuid.String(),
				Name: client.irc.Username,
				IP:   client.conn.RemoteAddr().String(),
			}

			result = append(result, details)
		}

		json.NewEncoder(w).Encode(result)
	}
}

func (server *server) serverListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(server.repository.servers)
	}
}

func (server *server) libraryHandler() http.HandlerFunc {
	type download struct {
		Name         string    `json:"name"`
		DownloadLink string    `json:"downloadLink"`
		Time         time.Time `json:"time"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		libraryDir := path.Join(server.config.DownloadDir, "books")
		books, err := os.ReadDir(libraryDir)
		if err != nil {
			server.log.Printf("Unable to list books. %s\n", err)
		}

		output := make([]download, 0)
		for _, book := range books {
			if !book.IsDir() {
				info, err := book.Info()
				if err != nil {
					server.log.Println(err)
				}

				dl := download{
					Name:         book.Name(),
					DownloadLink: path.Join("/static/library", book.Name()),
					Time:         info.ModTime(),
				}

				output = append(output, dl)
			}
		}

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	}
}

func (server *server) libraryFileHandler(route string) http.Handler {
	libraryPath := path.Join(server.config.DownloadDir, "books")
	fs := http.FileServer(http.Dir(libraryPath))
	return server.fileDeleter(http.StripPrefix(route, fs))
}

func (server *server) deleteLibraryFileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileName, err := url.PathUnescape(chi.URLParam(r, "fileName"))
		if err != nil {
			server.log.Printf("Error unescaping path: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		err = os.Remove(path.Join(server.config.DownloadDir, "books", fileName))
		if err != nil {
			server.log.Printf("Error deleting book file: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (server *server) testHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := server.getClient(r.Context())
		if client == nil {
			w.Write([]byte("Client not found."))
			return
		}
		fmt.Fprintf(w, "Your user id is %s", client.irc.Username)
	}
}
