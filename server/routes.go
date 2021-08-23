package server

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path"
)

//go:embed app/dist
var reactClient embed.FS

func registerRoutes(hub *Hub) {
	basePath := getBaseRoute()
	fmt.Printf("Base Path: %s\n", basePath)

	http.Handle(basePath, serveStaticFiles(basePath, "app/dist"))

	http.HandleFunc(path.Join(basePath, "ws"), func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	http.HandleFunc(path.Join(basePath, "connections"), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "There are currently %d active connections.", *numConnections)
	})
}

func serveStaticFiles(basePath, assetPath string) http.Handler {
	// update the embedded file system's tree so that index.html is at the root
	app, err := fs.Sub(reactClient, assetPath)
	if err != nil {
		log.Fatal(err)
	}

	// strip the predefined base path and serve the static file
	return http.StripPrefix(basePath, http.FileServer(http.FS(app)))
}

func getBaseRoute() string {
	cleaned := path.Clean(config.Basepath)
	if cleaned == "/" {
		return cleaned
	}
	return cleaned + "/"
}
