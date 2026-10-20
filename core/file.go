package core

import (
	"io"
	"os"
	"path/filepath"

	"github.com/evan-buss/openbooks/dcc"
	"github.com/evan-buss/openbooks/util"
)

func DownloadExtractDCCString(baseDir, dccStr string, progress io.Writer) (string, error) {
	// Download the file and wait until it is completed
	download, err := dcc.ParseString(dccStr)
	if err != nil {
		return "", err
	}

	dccPath := filepath.Join(baseDir, download.Filename+".temp")
	file, err := os.Create(dccPath)
	if err != nil {
		return "", err
	}

	writer := io.Writer(file)
	if progress != nil {
		writer = io.MultiWriter(file, progress)
	}

	// Download DCC data to the file
	if err := download.Download(writer); err != nil {
		// Close the handle before cleaning up; an open handle keeps the partial
		// file locked on Windows.
		file.Close()
		os.Remove(dccPath)
		return "", err
	}
	if err := file.Close(); err != nil {
		return "", err
	}
	if !util.IsArchive(dccPath) {
		return renameTempFile(dccPath), nil
	}

	extractedPath, err := util.ExtractArchive(dccPath)
	if err != nil {
		return "", err
	}

	return renameTempFile(extractedPath), nil
}

func renameTempFile(filePath string) string {
	if filepath.Ext(filePath) == ".temp" {
		newPath := filePath[:len(filePath)-len(".temp")]
		os.Rename(filePath, newPath)
		return newPath
	}

	return filePath
}
