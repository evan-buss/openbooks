package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/evan-buss/openbooks/core"
	"github.com/evan-buss/openbooks/irc"
)

type Config struct {
	UserName string // Username to use when connecting to IRC
	Log      bool   // True if IRC messages should be logged
	Dir      string
	Server   string
	irc      *irc.Conn
}

// StartInteractive instantiates the OpenBooks CLI interface
func StartInteractive(config Config) {
	fmt.Println("=======================================")
	fmt.Println("          Welcome to OpenBooks         ")
	fmt.Println("=======================================")

	instantiate(&config)
	defer config.irc.Close()

	ctx, cancel := context.WithCancel(context.Background())
	registerShutdown(config.irc, cancel)

	handler := fullHandler(config)
	if config.Log {
		file := config.setupLogger(handler)
		defer file.Close()
	}

	go core.StartReader(ctx, config.irc, handler)
	terminalMenu(config.irc)

	<-ctx.Done()
}

func StartDownload(config Config, download string) {
	instantiate(&config)
	defer config.irc.Close()
	ctx, cancel := context.WithCancel(context.Background())

	handler := core.EventHandler{}
	addRequiredHandlers(handler, &config)
	handler[core.BookResult] = func(text string) {
		fmt.Printf("%sReceived file response.\n", clearLine)
		config.downloadHandler(text)
		cancel()
	}
	if config.Log {
		file := config.setupLogger(handler)
		defer file.Close()
	}

	fmt.Printf("Sending download request.")
	go core.StartReader(ctx, config.irc, handler)
	core.DownloadBook(config.irc, download)
	fmt.Printf("%sSent download request.", clearLine)
	fmt.Printf("Waiting for file response.")

	registerShutdown(config.irc, cancel)
	<-ctx.Done()
}

func StartSearch(config Config, query string) {
	nextSearchTime := GetLastSearchTime().Add(15 * time.Second)
	instantiate(&config)
	defer config.irc.Close()
	ctx, cancel := context.WithCancel(context.Background())

	handler := core.EventHandler{}
	addRequiredHandlers(handler, &config)
	handler[core.SearchResult] = func(text string) {
		fmt.Printf("%sReceived file response.\n", clearLine)
		config.searchHandler(text)
		cancel()
	}
	handler[core.MatchesFound] = config.matchesFoundHandler
	if config.Log {
		file := config.setupLogger(handler)
		defer file.Close()
	}

	fmt.Printf("Sending search request.")
	time.Sleep(time.Until(nextSearchTime))
	go core.StartReader(ctx, config.irc, handler)
	core.SearchBook(config.irc, query)
	SetLastSearchTime()
	fmt.Printf("%sSent search request.", clearLine)
	fmt.Printf("Waiting for file response.")

	registerShutdown(config.irc, cancel)
	<-ctx.Done()
}
