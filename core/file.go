package core

import (
	"fmt"
	"os"
	"path"

	"github.com/evan-buss/openbooks/dcc"
	"github.com/evan-buss/openbooks/util"
)

func DownloadExtractDCCString(baseDir, dccStr string, progress dcc.ProgressFunc) (string, error) {
	// Download the file and wait until it is completed
	download, err := dcc.ParseString(dccStr)
	if err != nil {
		return "", err
	}

	dccPath := path.Join(baseDir, download.Filename)
	file, err := os.Create(dccPath)
	if err != nil {
		return "", err
	}
	fmt.Println()
	// Download DCC data to the file
	err = download.Download(file, progress)
	fmt.Println()
	if err != nil {
		return "", err
	}
	file.Close()
	if !util.IsArchive(dccPath) {
		return dccPath, nil
	}

	extractedPath, err := util.ExtractArchive(dccPath)
	if err != nil {
		return "", err
	}

	return extractedPath, nil
}
