package search

import (
	"bufio"
	"context"
	"github.com/evan-buss/openbooks/core"
	"go.uber.org/zap"
	"log"
	"os"
	"path/filepath"
	"strings"

	openbooksv1 "github.com/evan-buss/openbooks/proto/gen/proto/go/openbooks/v1"
)

type Indexer interface {
	Index(ctx context.Context) <-chan *openbooksv1.Book
}

// DirectoryIndexer reads all text files from a specific folder
// and parses the lines in each file to *Book
type DirectoryIndexer struct {
	Path string
	log  *zap.SugaredLogger
}

func NewDirectoryIndexer(path string, log *zap.SugaredLogger) *DirectoryIndexer {
	return &DirectoryIndexer{path, log}
}

// Index reads all text files from directory and indexes them.
func (i *DirectoryIndexer) Index(ctx context.Context) <-chan *openbooksv1.Book {
	out := make(chan *openbooksv1.Book)
	go func() {
		files, err := os.ReadDir(i.Path)
		if err != nil {
			log.Fatal(err)
		}
		for _, dirEntry := range files {
			if ctx.Err() != nil {
				return
			}

			if dirEntry.IsDir() || filepath.Ext(dirEntry.Name()) != ".txt" {
				continue
			}

			i.log.Infof("indexing file %s", dirEntry.Name())

			f, err := os.Open(filepath.Join(i.Path, dirEntry.Name()))
			if err != nil {
				log.Println(err)
				continue
			}

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "!") {
					detail, _ := core.ParseLineV2(line)
					out <- &openbooksv1.Book{
						Server: detail.Server,
						Author: detail.Author,
						Title:  detail.Title,
						Format: detail.Format,
						Size:   detail.Size,
						Full:   detail.Full,
						Line:   line,
					}
				}
			}
			err = f.Close()
			if err != nil {
				i.log.Errorw("unable to close file", "err", err)
			}
		}

		close(out)
	}()

	return out
}
