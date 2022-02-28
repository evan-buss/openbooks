package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"

	"github.com/evan-buss/openbooks/core"
)

type Indexer interface {
	Index() <-chan Book
}

type DirectoryIndexer struct {
	Path string
}

// Index reads all text files from directory and indexes them.
func (i *DirectoryIndexer) Index() <-chan Book {
	out := make(chan Book)
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

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				detail, err := core.ParseLine(scanner.Text())
				if err != nil {
					// log.Println(err)
				} else {
					out <- NewBook(detail, scanner.Text())
				}
				// log.Println(detail)
			}
			f.Close()
		}

		close(out)
	}()

	return out
}
