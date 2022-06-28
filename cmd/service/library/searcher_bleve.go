package library

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/v2/analysis/lang/en"
	"github.com/blevesearch/bleve/v2/mapping"
	openbooksv1 "github.com/evan-buss/openbooks/proto/gen/proto/go/openbooks/v1"
	"google.golang.org/protobuf/types/known/durationpb"
)

type BleveSearcher struct {
	Index bleve.Index
}

func NewBleveSeacher(indexFolder string) (*BleveSearcher, error) {
	bookMapping := bleve.NewDocumentMapping()

	serverMapping := bleve.NewKeywordFieldMapping()
	serverMapping.Analyzer = keyword.Name
	serverMapping.Store = false
	bookMapping.AddFieldMappingsAt("server", serverMapping)

	authorMapping := bleve.NewTextFieldMapping()
	authorMapping.Analyzer = en.AnalyzerName
	authorMapping.Store = false
	bookMapping.AddFieldMappingsAt("author", authorMapping)

	titleMapping := bleve.NewTextFieldMapping()
	titleMapping.Analyzer = en.AnalyzerName
	titleMapping.Store = false
	bookMapping.AddFieldMappingsAt("title", titleMapping)

	formatMapping := bleve.NewKeywordFieldMapping()
	formatMapping.Analyzer = keyword.Name
	formatMapping.Store = false
	bookMapping.AddFieldMappingsAt("format", formatMapping)

	lineMapping := bleve.NewKeywordFieldMapping()
	lineMapping.Analyzer = en.AnalyzerName
	lineMapping.Store = true
	bookMapping.AddFieldMappingsAt("line", lineMapping)

	sizeMapping := bleve.NewDocumentDisabledMapping()
	bookMapping.AddSubDocumentMapping("size", sizeMapping)

	fullMapping := bleve.NewDocumentDisabledMapping()
	bookMapping.AddSubDocumentMapping("full", fullMapping)

	mapping := bleve.NewIndexMapping()
	mapping.DefaultType = "book"
	mapping.AddDocumentMapping("book", bookMapping)

	index, err := useIndexFolder(indexFolder, mapping)
	if err != nil {
		return nil, err
	}

	return &BleveSearcher{Index: index}, nil
}

func (b *BleveSearcher) Search(ctx context.Context, q string) (*openbooksv1.SearchResponse, error) {
	stats, _ := b.Index.Stats().MarshalJSON()
	log.Println(string(stats))

	// search for some text
	query := bleve.NewQueryStringQuery(q)
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"line"}
	search.Size = 1000
	searchResults, err := b.Index.SearchInContext(ctx, search)
	if err != nil {
		return nil, err
	}
	results := make([]string, 0, searchResults.Total)

	for _, hit := range searchResults.Hits {
		line, ok := hit.Fields["line"].(string)
		if !ok {
			return nil, err
		}
		results = append(results, line)
	}

	response := &openbooksv1.SearchResponse{
		Total:    searchResults.Total,
		Results:  results,
		Duration: durationpb.New(searchResults.Took),
	}

	return response, nil
}

func (b *BleveSearcher) AddDocuments(ctx context.Context, documents <-chan *openbooksv1.Book) (time.Duration, error) {
	batchSize := 1000
	count := 0
	startTime := time.Now()
	batch := b.Index.NewBatch()
	batchCount := 0

	for bookDetail := range documents {
		hash := md5.Sum([]byte(bookDetail.Full))
		hashStr := hex.EncodeToString(hash[:])

		if ctx.Err() != nil {
			return 0, ctx.Err()
		}

		// Don't index if we already have the line in the index
		doc, err := b.Index.Document(hashStr)
		if doc != nil {
			continue
		}

		if err != nil {
			return 0, ctx.Err()
		}

		err = batch.Index(hashStr, bookDetail)
		if err != nil {
			return 0, err
		}

		batchCount++

		if batchCount >= batchSize {
			err := b.Index.Batch(batch)
			if err != nil {
				return 0, err
			}
			batch.Reset()
			batchCount = 0
		}
		count++
	}
	// flush the last batch
	if batchCount > 0 {
		err := b.Index.Batch(batch)
		if err != nil {
			return 0, fmt.Errorf("flushing last batch of %d documents. %w", batchCount, err)
		}
	}

	return time.Since(startTime), nil
}

func useIndexFolder(path string, mapping *mapping.IndexMappingImpl) (bleve.Index, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		index, err := bleve.New(path, mapping)
		if err != nil {
			return nil, err
		}
		return index, err
	}

	return bleve.Open(path)
}
