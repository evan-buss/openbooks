package bleve

import (
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

func NewBleveSeacher() (*BleveSearcher, error) {
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

	index, err := useIndexFolder("openbooks.bleve", mapping)
	if err != nil {
		return nil, err
	}

	return &BleveSearcher{Index: index}, nil
}

func (b *BleveSearcher) Search(q string) openbooksv1.SearchResponse {
	// search for some text
	query := bleve.NewMatchQuery(q)
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"line"}
	search.Size = 1000
	searchResults, err := b.Index.Search(search)
	if err != nil {
		fmt.Println(err)
		return openbooksv1.SearchResponse{
			Duration: durationpb.New(0),
			Results:  []string{},
			Total:    0,
		}
	}
	results := make([]string, 0, searchResults.Total)

	for _, hit := range searchResults.Hits {
		if line, ok := hit.Fields["line"].(string); ok {
			results = append(results, line)
			continue
		}
		log.Println("Unable to get field 'Line' from search result. Check SearchRequest.Fields...")
	}

	return openbooksv1.SearchResponse{
		Total:    searchResults.Total,
		Results:  results,
		Duration: durationpb.New(searchResults.Took),
	}
}

func (b *BleveSearcher) AddDocument(doc *openbooksv1.Book) {
	hash := md5.Sum([]byte(doc.Full))
	err := b.Index.Index(hex.EncodeToString(hash[:]), doc)
	if err != nil {
		log.Println(err)
	}
}

func (b *BleveSearcher) AddDocuments(documents <-chan *openbooksv1.Book) time.Duration {
	batchSize := 1000
	count := 0
	startTime := time.Now()
	batch := b.Index.NewBatch()
	batchCount := 0
	for bookDetail := range documents {
		hash := md5.Sum([]byte(bookDetail.Full))
		batch.Index(hex.EncodeToString(hash[:]), bookDetail)

		batchCount++

		if batchCount >= batchSize {
			err := b.Index.Batch(batch)
			if err != nil {
				log.Fatal(err)
			}
			batch.Reset()
			batchCount = 0
		}
		count++
		if count%1000 == 0 {
			indexDuration := time.Since(startTime)
			indexDurationSeconds := float64(indexDuration) / float64(time.Second)
			timePerDoc := float64(indexDuration) / float64(count)
			log.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
		}
	}
	// flush the last batch
	if batchCount > 0 {
		err := b.Index.Batch(batch)
		if err != nil {
			log.Fatal(err)
		}
	}
	indexDuration := time.Since(startTime)
	indexDurationSeconds := float64(indexDuration) / float64(time.Second)
	timePerDoc := float64(indexDuration) / float64(count)
	log.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))

	return indexDuration
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
