package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/evan-buss/openbooks/core"
	"github.com/evan-buss/openbooks/irc"
	"github.com/evan-buss/openbooks/util"
)

const clearLine = "\r\033[2K"

func registerShutdown(conn *irc.Conn, cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		conn.Disconnect()
		cancel()
		os.Exit(0)
	}()
}

// Connect to IRC server and save connection to Config
func instantiate(config *Config) {
	fmt.Printf("Connecting to %s.", config.Server)
	conn := irc.New(config.UserName, "OpenBooks CLI")
	config.irc = conn
	core.Join(conn, config.Server)
	fmt.Printf("%sConnected to %s.\n", clearLine, config.Server)
}

// Required handlers are used regardless of what CLI mode is selected.
// Keep alive pings and other core IRC client features
func addRequiredHandlers(handler core.EventHandler, config *Config) {
	handler[core.Ping] = config.pingHandler
	handler[core.Version] = config.versionHandler
}

func (config *Config) setupLogger(handler core.EventHandler) io.Closer {
	logger, file, err := util.CreateLogFile(config.UserName, config.Dir)
	if err != nil {
		log.Fatalf("Error setting up logger: %s\n", err)
	}
	handler[core.Message] = func(text string) {
		logger.Println(text)
	}

	return file
}

func GetLastSearchTime() time.Time {
	timestampFilePath := filepath.Join(os.TempDir(), ".openbooks")
	fileInfo, err := os.Stat(timestampFilePath)

	if errors.Is(err, os.ErrNotExist) {
		return time.Time{}
	}

	return fileInfo.ModTime()
}

func SetLastSearchTime() {
	timestampFilePath := filepath.Join(os.TempDir(), ".openbooks")
	_, err := os.Stat(timestampFilePath)

	if errors.Is(err, os.ErrNotExist) {
		os.Create(timestampFilePath)
	}

	os.Chtimes(timestampFilePath, time.Now(), time.Now())
}
