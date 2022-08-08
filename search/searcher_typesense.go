package search

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	openbooksv1 "github.com/evan-buss/openbooks/proto/gen/proto/go/openbooks/v1"
	"github.com/samber/lo"
	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
	"google.golang.org/protobuf/types/known/durationpb"
)

type TypesenseSearcher struct {
	client *typesense.Client
}

func NewTypesenseSearcher(url, masterKey string) (*TypesenseSearcher, error) {
	client := typesense.NewClient(
		typesense.WithServer(url),
		typesense.WithAPIKey(masterKey),
	)

	return &TypesenseSearcher{client}, nil
}

func Ptr[T any](v T) *T {
	return &v
}

func (s *TypesenseSearcher) Search(ctx context.Context, query string) (*openbooksv1.SearchResponse, error) {
	start := time.Now()
	result, err := s.client.Collection("books").Documents().Search(&api.SearchCollectionParams{
		Q:       query,
		QueryBy: "title, author",
		PerPage: Ptr(250),
	})

	return &openbooksv1.SearchResponse{
		Results: lo.Map(*result.Hits, func(hit api.SearchResultHit, i int) string {
			return (*hit.Document)["full"].(string)
		}),
		Duration: durationpb.New(time.Since(start)),
		Total:    uint64(*result.Found),
	}, err
}

func (s *TypesenseSearcher) AddDocuments(ctx context.Context, books <-chan *openbooksv1.Book) (time.Duration, error) {
	fmt.Println("start typesense indexing")
	startTime := time.Now()

	batchSize := 1000

	batch := make([]interface{}, batchSize)

	count := 0
	errors := 0

	index := 0
	for bookDetail := range books {
		if index < batchSize {
			hash := md5.Sum([]byte(bookDetail.Full))
			hashStr := hex.EncodeToString(hash[:])
			batch[index] = BookWrapper{hashStr, bookDetail}
			index++
		} else {
			fmt.Println("batch size reached, sending batch")
			batch = batch[:index]
			index = 0
			_, err := s.client.Collection("books").Documents().Import(batch, &api.ImportDocumentsParams{})
			fmt.Println("batch sent")
			if err != nil {
				errors++
			}
			batch = make([]interface{}, batchSize)
		}
	}

	fmt.Println("done typesense indexing")
	fmt.Println("errors: ", errors)
	fmt.Println("count: ", count)
	return time.Since(startTime), nil
}
