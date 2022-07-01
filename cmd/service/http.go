package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (s *server) SwaggerDocHandlers() *chi.Mux {
	router := chi.NewRouter()

	// Register Swagger Documentation
	router.Get("/", http.RedirectHandler("/swagger/index.html", http.StatusTemporaryRedirect).ServeHTTP)
	router.Get("/openbooks.swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../proto/gen/docs/openbooks/v1/openbooks.swagger.json")
	})
	router.Get("/*", httpSwagger.Handler(
		httpSwagger.URL("./openbooks.swagger.json"),
	))

	return router
}
