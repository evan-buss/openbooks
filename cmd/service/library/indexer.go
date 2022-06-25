package library

import (
	"log"
	"os"
	"path/filepath"

	"github.com/evan-buss/openbooks/core"
	openbooksv1 "github.com/evan-buss/openbooks/proto/gen/proto/go/openbooks/v1"
)

// func StartIndex(directory string) error {
// 	entries, err := os.ReadDir(directory)
// 	if err != nil {
// 		return err
// 	}

// 	for _, entry := range entries {
// 		if entry.IsDir() || filepath.Ext(entry.Name()) != ".zip" {
// 			continue
// 		}

// 		archiveDir := filepath.Join(directory, entry.Name())
// 		filename := entry.Name()[:len(entry.Name())-len(filepath.Ext(entry.Name()))] + ".txt"
// 		textFilePath := filepath.Join(directory, filename+".txt")
// 		log.Println(archiveDir, filename, textFilePath)
// 		err = archiver.Extract(archiveDir, filename, textFilePath)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

type Indexer interface {
	Index() <-chan Book
}

type DirectoryIndexer struct {
	Path string
}

// Index reads all text files from directory and indexes them.
func (i *DirectoryIndexer) Index() <-chan *openbooksv1.Book {
	out := make(chan *openbooksv1.Book)
	go func() {
		files, err := os.ReadDir(i.Path)
		if err != nil {
			log.Fatal(err)
		}
		for _, dirEntry := range files {
			if dirEntry.IsDir() || filepath.Ext(dirEntry.Name()) != ".txt" {
				continue
			}

			log.Println("Indexing File: ", dirEntry.Name())

			f, err := os.Open(filepath.Join(i.Path, dirEntry.Name()))
			if err != nil {
				log.Println(err)
				continue
			}

			books, _ := core.ParseSearchV2(f)
			for _, detail := range books {
				out <- &openbooksv1.Book{
					Server: detail.Server,
					Author: detail.Author,
					Title:  detail.Title,
					Format: detail.Format,
					Size:   detail.Size,
					Full:   detail.Full,
				}
			}

			// scanner := bufio.NewScanner(f)
			// for scanner.Scan() {
			// 	detail, err := core.ParseLineV2(scanner.Text())
			// 	if err != nil {
			// 		// log.Println(err)
			// 	} else {
			// 		out <- &openbooksv1.Book{
			// 			Server: detail.Server,
			// 			Author: detail.Author,
			// 			Title:  detail.Title,
			// 			Format: detail.Format,
			// 			Size:   detail.Size,
			// 			Full:   detail.Full,
			// 			Line:   scanner.Text(),
			// 		}
			// 	}
			// 	// log.Println(detail)
			// }
			f.Close()
		}

		close(out)
	}()

	return out
}
