package util

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"
)

func CreateLogFile(username, dir string) (*log.Logger, io.Closer, error) {
	date := time.Now().Format("2006-01-02--15-04-05")
	fileName := fmt.Sprintf("%s--%s.log", username, date)

	err := os.MkdirAll(path.Join(dir, "logs"), os.FileMode(0755))
	if err != nil {
		return nil, nil, err
	}

	path := path.Join(dir, "logs", fileName)
	logFile, err := os.Create(path)
	if err != nil {
		return nil, nil, err
	}

	return log.New(logFile, "", 0), logFile, nil
}
