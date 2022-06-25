package library

import (
	"time"

	openbooksv1 "github.com/evan-buss/openbooks/proto/gen/proto/go/openbooks/v1"
)

type Searcher interface {
	Search(query string) openbooksv1.SearchResponse
	AddDocument(book *openbooksv1.Book)
	AddDocuments(<-chan *openbooksv1.Book) time.Duration
}
