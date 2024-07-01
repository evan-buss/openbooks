package server

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/evan-buss/openbooks/core"

	"github.com/go-chi/chi/v5"
)

//go:embed app/dist
var reactClient embed.FS

func (server *server) registerRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Handle("/*", server.staticFilesHandler("app/dist"))
	router.HandleFunc("/events", server.eventsHandler())
	// router.Get("/stats", server.statsHandler())
	router.Get("/servers", server.serverListHandler())

	router.Group(func(r chi.Router) {
		r.Post("/search", server.searchHandler())
		r.Post("/download", server.downloadHandler())

		r.Get("/library", server.getAllBooksHandler())
		r.Delete("/library/{fileName}", server.deleteBooksHandler())
		r.Get("/library/*", server.getBookHandler())
	})

	return router
}

func (server *server) eventsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if server.connectedClients.Load() > 0 {
			http.Error(w, "Multiple connections not allowed.", http.StatusServiceUnavailable)
			return
		}

		server.log.Printf("Client connected from %s\n", r.RemoteAddr)

		// The SSE library expects a "stream" query parameter to be set
		newQuery := r.URL.Query()
		newQuery.Set("stream", "events")
		r.URL.RawQuery = newQuery.Encode()

		go func() {
			<-r.Context().Done()
			server.log.Println("Client disconnected from", r.RemoteAddr)
		}()

		server.sse.ServeHTTP(w, r)
	}
}

func (server *server) staticFilesHandler(assetPath string) http.Handler {
	// update the embedded file system's tree so that index.html is at the root
	app, err := fs.Sub(reactClient, assetPath)
	if err != nil {
		server.log.Println(err)
	}

	// strip the predefined base path and serve the static file
	return http.StripPrefix(server.config.Basepath, http.FileServer(http.FS(app)))
}

// func (server *server) statsHandler() http.HandlerFunc {
// 	type statsReponse struct {
// 		UUID string `json:"uuid"`
// 		Name string `json:"name"`
// 	}

// 	return func(w http.ResponseWriter, r *http.Request) {
// 		result := make([]statsReponse, 0, len(server.clients))

// 		for _, client := range server.clients {
// 			details := statsReponse{
// 				UUID: client.uuid.String(),
// 				Name: client.irc.Username,
// 			}

// 			result = append(result, details)
// 		}

// 		json.NewEncoder(w).Encode(result)
// 	}
// }

func (server *server) serverListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(server.repository.servers)
	}
}

func (server *server) getAllBooksHandler() http.HandlerFunc {
	type download struct {
		Name         string    `json:"name"`
		DownloadLink string    `json:"downloadLink"`
		Time         time.Time `json:"time"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if !server.config.Persist {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		libraryDir := filepath.Join(server.config.DownloadDir, "books")
		books, err := os.ReadDir(libraryDir)
		if err != nil {
			server.log.Printf("Unable to list books. %s\n", err)
		}

		output := make([]download, 0)
		for _, book := range books {
			if book.IsDir() || strings.HasPrefix(book.Name(), ".") || filepath.Ext(book.Name()) == ".temp" {
				continue
			}

			info, err := book.Info()
			if err != nil {
				server.log.Println(err)
			}

			dl := download{
				Name:         book.Name(),
				DownloadLink: path.Join("library", book.Name()),
				Time:         info.ModTime(),
			}

			output = append(output, dl)
		}

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	}
}

func (server *server) getBookHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, fileName := path.Split(r.URL.Path)
		bookPath := filepath.Join(server.config.DownloadDir, "books", fileName)

		http.ServeFile(w, r, bookPath)

		if !server.config.Persist {
			err := os.Remove(bookPath)
			if err != nil {
				server.log.Printf("Error when deleting book file. %s", err)
			}
		}
	}
}

func (server *server) deleteBooksHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileName, err := url.PathUnescape(chi.URLParam(r, "fileName"))
		if err != nil {
			server.log.Printf("Error unescaping path: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		err = os.Remove(filepath.Join(server.config.DownloadDir, "books", fileName))
		if err != nil {
			server.log.Printf("Error deleting book file: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (server *server) searchHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		server.lastSearchMutex.Lock()
		defer server.lastSearchMutex.Unlock()

		nextAvailableSearch := server.lastSearch.Add(server.config.SearchTimeout)

		if time.Now().Before(nextAvailableSearch) {
			remainingSeconds := time.Until(nextAvailableSearch).Seconds()
			// TODO: Show HTTP errors on client instead of sending SSE message
			http.Error(w, "Rate limited.", http.StatusTooManyRequests)
			server.send <- newRateLimitResponse(remainingSeconds)
			return
		}

		query := r.URL.Query().Get("query")
		if query == "" {
			http.Error(w, "No search query provided.", http.StatusBadRequest)
			server.send <- newErrorResponse("No search query provided.")
			return
		}

		core.SearchBook(server.irc, server.config.SearchBot, query)
		server.lastSearch = time.Now()

		w.WriteHeader(http.StatusAccepted)

		server.send <- newStatusResponse(NOTIFY, "Search request sent.")
	}
}

func (server *server) downloadHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		book := r.URL.Query().Get("book")
		if book == "" {
			http.Error(w, "No book provided.", http.StatusBadRequest)
			server.send <- newErrorResponse("No book provided.")
			return
		}

		core.DownloadBook(server.irc, book)
		w.WriteHeader(http.StatusAccepted)

		server.send <- newStatusResponse(NOTIFY, "Download request received.")
	}
}
