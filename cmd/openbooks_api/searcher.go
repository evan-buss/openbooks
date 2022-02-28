package main

import (
	"time"
)

type SearchResultCollection struct {
	Results  []string `json:"results"`
	Duration int64    `json:"duration"`
	Total    uint64   `json:"total"`
}

type Searcher interface {
	Search(string) SearchResultCollection
	AddDocument(Book)
	AddDocuments(<-chan Book) time.Duration
}
