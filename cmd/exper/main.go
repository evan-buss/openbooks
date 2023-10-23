package main

import (
	"net/http"
	"os"

	"log/slog"

	"github.com/evan-buss/openbooks/opds"
	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()
	opdsRouter := opds.NewOPDS("C:/Users/evanb/Downloads/Books", "C:/Users/evanb/Downloads/Books/SearchResults")
	router.Mount("/opds", opdsRouter)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Default().WithGroup("main").Info("Starting server")

	http.ListenAndServe(":8080", router)
}
