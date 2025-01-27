package util

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/mholt/archiver/v3"
)

var (
	ErrNotFullyCopied = errors.New("didn't copy entire file from the archive")
)

func ExtractArchive(archivePath string) (string, error) {
	// Our path will have a .temp appended to it so we can't rely on the automatic file-extension based archive extractor selection.
	// This code was taken from the archiver.Walk(archive string, walkFn WalkFunc) error function.
	// We just remove .temp before trying to find a matching archive extractor.
	wIface, err := archiver.ByExtension(archivePath[:len(archivePath)-len(".temp")])
	if err != nil {
		return "", err
	}
	w, ok := wIface.(archiver.Walker)
	if !ok {
		return "", fmt.Errorf("format specified by archive filename is not a walker format: %s (%T)", archivePath, wIface)
	}

	var newPath string
	err = w.Walk(archivePath, func(f archiver.File) error {
		// Extract only one file per archive. Otherwise, stop walking,
		// remove extracted items, and deliver the archive itself.
		if newPath != "" {
			err := os.Remove(newPath)
			if err != nil {
				return err
			}
			newPath = ""
			return archiver.ErrStopWalk
		}

		newPath = filepath.Join(filepath.Dir(archivePath), f.Name()+".temp")

		out, err := os.Create(newPath)
		if err != nil {
			return err
		}

		copied, err := io.Copy(out, f)
		if err != nil {
			return err
		}
		if copied != f.Size() {
			return ErrNotFullyCopied
		}

		err = out.Close()
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	// If we extracted exactly one file, send that file and remove the zip file.
	// Otherwise, send the archive itself.
	if newPath != "" {
		err := os.Remove(archivePath)
		if err != nil {
			log.Println("remove error", err)
		}
		return newPath, nil
	} else {
		return archivePath, nil
	}
}

// IsArchive returns true if the file at the given path is an archive that can
// be extracted. Returns false otherwise.
func IsArchive(path string) bool {
	if filepath.Ext(path) == ".temp" {
		path = path[:len(path)-len(".temp")]
	}

	_, err := archiver.ByExtension(path)
	return err == nil
}
