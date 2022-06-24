package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (s *server) SwaggerDocHandlers() *chi.Mux {
	router := chi.NewRouter()

	url := fmt.Sprintf("http://localhost:%s/swagger", s.config.port)
	// Register Swagger Documentation
	router.Get("/", http.RedirectHandler(url+"/index.html", http.StatusTemporaryRedirect).ServeHTTP)
	router.Get("/openbooks.swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../proto/gen/docs/openbooks/v1/openbooks.swagger.json")
	})
	router.Get("/*", httpSwagger.Handler(
		httpSwagger.URL(url+"/openbooks.swagger.json"),
	))

	return router
}
