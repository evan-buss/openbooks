package library

import (
	"context"
	"time"

	openbooksv1 "github.com/evan-buss/openbooks/proto/gen/proto/go/openbooks/v1"
)

type Searcher interface {
	Search(ctx context.Context, query string) (*openbooksv1.SearchResponse, error)
	AddDocuments(ctx context.Context, books <-chan *openbooksv1.Book) (time.Duration, error)
}
