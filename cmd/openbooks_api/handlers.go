package main

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
)

func (s *server) registerRoutes() {
	s.router.Get("/search", s.limit(s.searchHandler()).ServeHTTP)
	s.router.Post("/index", s.indexHandler())
	s.router.Get("/stats", s.statsHandler())
}

func (s *server) searchHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		results := s.searcher.Search(query)
		json.NewEncoder(w).Encode(results)
	}
}

func (s *server) indexHandler() http.HandlerFunc {
	var indexLock uint32

	return func(w http.ResponseWriter, r *http.Request) {
		if !atomic.CompareAndSwapUint32(&indexLock, 0, 1) {
			w.Write([]byte("Index already in progress..."))
			return
		}
		defer atomic.StoreUint32(&indexLock, 0)

		s.searcher.AddDocuments(s.indexer.Index())
	}
}

func (s *server) statsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bleve := s.searcher.(*BleveSearcher)
		stats, _ := bleve.Index.Stats().MarshalJSON()
		w.Write(stats)
	}
}
