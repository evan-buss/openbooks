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
	"strings"
	"syscall"
	"time"

	"github.com/evan-buss/openbooks/core"
	"github.com/evan-buss/openbooks/util"
	"github.com/schollz/progressbar/v3"
)

var servers []string

const clearLine = "\r\033[2K"

func registerShutdown(client *core.IrcClient, cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		client.Disconnect()
		cancel()
		os.Exit(0)
	}()
}

// Connect to IRC server and save the client to Config
func instantiate(config *Config) {
	fmt.Printf("Connecting to %s.", config.Server)
	config.client = core.NewIrcClient(config.UserName, config.Version, config.Dir)
	err := config.client.Connect(config.Server, config.EnableTLS)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%sConnected to %s.\n", clearLine, config.Server)
}

// baseHandlers returns the Handlers common to every CLI mode: server-list
// tracking, a terminal progress bar, and (optionally) raw IRC logging. The
// returned io.Closer (the log file) is nil when logging is disabled.
func (config *Config) baseHandlers() (core.Handlers, io.Closer) {
	handlers := core.Handlers{
		ServerList: func(s core.IrcServers) { servers = s.ElevatedUsers },
		Progress: func(filename string, size int64) io.Writer {
			return progressbar.DefaultBytes(size, filename)
		},
		Erred: func(err error) { log.Println(err) },
	}

	if !config.Log {
		return handlers, nil
	}

	logger, file, err := util.CreateLogFile(config.UserName, config.Dir)
	if err != nil {
		log.Fatalf("Error setting up logger: %s\n", err)
	}
	handlers.Message = func(text string) { logger.Println(text) }
	return handlers, file
}

// Show warning message if the server they are downloading from is not online.
func warnIfServerOffline(bookLine string) {
	for _, server := range servers {
		if strings.HasPrefix(bookLine[1:], server) {
			return
		}
	}

	fmt.Println("WARNING: That server is not online. Your request will never complete.")
}

func getLastSearchTime() time.Time {
	timestampFilePath := filepath.Join(os.TempDir(), ".openbooks")
	fileInfo, err := os.Stat(timestampFilePath)

	if errors.Is(err, os.ErrNotExist) {
		return time.Now()
	}

	return fileInfo.ModTime()
}

func setLastSearchTime() {
	timestampFilePath := filepath.Join(os.TempDir(), ".openbooks")
	_, err := os.Stat(timestampFilePath)

	if errors.Is(err, os.ErrNotExist) {
		os.Create(timestampFilePath)
	}

	os.Chtimes(timestampFilePath, time.Now(), time.Now())
}
