package search

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	openbooksv1 "github.com/evan-buss/openbooks/proto/gen/proto/go/openbooks/v1"
	"github.com/meilisearch/meilisearch-go"
)

type MeiliSearcher struct {
	client *meilisearch.Client
}

func NewMeiliSearcher() (*MeiliSearcher, error) {
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   "http://localhost:7700",
		APIKey: "masterKey",
	})

	return &MeiliSearcher{client}, nil
}

func (s *MeiliSearcher) Search(ctx context.Context, query string) (*openbooksv1.SearchResponse, error) {
	panic("not implemented") // TODO: Implement
}

type BookWrapper struct {
	Id string `json:"id"`
	*openbooksv1.Book
}

func (s *MeiliSearcher) AddDocuments(ctx context.Context, books <-chan *openbooksv1.Book) (time.Duration, error) {

	mIndex, err := s.client.GetIndex("books")
	if err != nil {
		s.client.CreateIndex(&meilisearch.IndexConfig{
			Uid:        "books",
			PrimaryKey: "id",
		})

		mIndex = s.client.Index("books")
	}

	fmt.Println("start meilisearch indexing")
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
			fmt.Println("SENDING BATCH")
			batch = batch[:index]
			index = 0
			task, err := mIndex.AddDocuments(batch)
			if err != nil {
				fmt.Println(err)
				errors++
			}
			_, err = s.client.WaitForTask(task.TaskUID, meilisearch.WaitParams{Context: ctx, Interval: time.Millisecond * 100})
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("BATCH SENT")
			batch = make([]interface{}, batchSize)
		}
	}

	fmt.Println("done typesense indexing")
	fmt.Println("errors: ", errors)
	fmt.Println("count: ", count)
	return time.Since(startTime), nil
}
