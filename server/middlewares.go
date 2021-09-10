package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/google/uuid"
)

type userCtxKeyType string
type uuidCtxKeyType string

const (
	userCtxKey userCtxKeyType = "user-client"
	uuidCtxKey uuidCtxKeyType = "user-uuid"
)

func (server *server) requireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("OpenBooks")
		if err != nil {
			server.log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userUUID, err := uuid.Parse(cookie.Value)
		if err != nil {
			server.log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), uuidCtxKey, userUUID)))
	})
}

// getClient should only be called when requireUser is in the middleware chain.
func (server *server) getClient(ctx context.Context) *Client {

	user := getUUID(ctx)
	if user == uuid.Nil {
		return nil
	}

	if client, ok := server.clients[user]; ok {
		return client
	}

	return nil
}

func getUUID(ctx context.Context) uuid.UUID {
	uid, ok := ctx.Value(uuidCtxKey).(uuid.UUID)
	if !ok {
		log.Println("Unable to find user.")
	}
	return uid
}

// fileDeleter handles deleting files after a download is successful.
func (server *server) fileDeleter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		if !server.config.Persist {
			_, fileName := path.Split(r.URL.Path)
			bookPath := path.Join(server.config.DownloadDir, "books", fileName)
			err := os.Remove(bookPath)
			if err != nil {
				server.log.Printf("Error when deleting book file. %s", err)
			}
		}
	})
}
