package util

import (
	"errors"
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
	var newPath string
	err := archiver.Walk(archivePath, func(f archiver.File) error {
		newPath = filepath.Join(filepath.Dir(archivePath), f.Name())

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

	// If we extracted a file, send that file and remove the zip file
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
	_, err := archiver.ByExtension(path)
	return err == nil
}
