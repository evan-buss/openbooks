package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/evan-buss/openbooks/core"
)

type Config struct {
	UserName  string // Username to use when connecting to IRC
	Log       bool   // True if IRC messages should be logged
	Dir       string
	Server    string
	EnableTLS bool
	SearchBot string
	Version   string
	client    *core.IrcClient
}

// StartInteractive instantiates the OpenBooks CLI interface
func StartInteractive(config Config) {
	fmt.Println("=======================================")
	fmt.Println("          Welcome to OpenBooks         ")
	fmt.Println("=======================================")

	instantiate(&config)
	defer config.client.Disconnect()

	ctx, cancel := context.WithCancel(context.Background())
	registerShutdown(config.client, cancel)

	handlers, closer := config.baseHandlers()
	if closer != nil {
		defer closer.Close()
	}
	handlers.SearchResults = func(_ []core.BookDetail, _ []core.ParseError, path string) {
		fmt.Println("Results location: " + path)
		terminalMenu(config)
	}
	handlers.BookDownloaded = func(path string) {
		fmt.Println("File location: " + path)
		terminalMenu(config)
	}
	handlers.NoResults = func() {
		fmt.Println("No results returned for your search...")
		terminalMenu(config)
	}
	handlers.BadServer = func() {
		fmt.Println("That server is not available. Try again...")
		terminalMenu(config)
	}
	handlers.SearchAccepted = func() { fmt.Println("Search has been accepted. Please wait.") }
	handlers.MatchesFound = func(num string) { fmt.Printf("Found %s search results.", num) }

	config.client.StartReader(ctx, handlers)
	terminalMenu(config)

	<-ctx.Done()
}

func StartDownload(config Config, download string) {
	instantiate(&config)
	defer config.client.Disconnect()
	ctx, cancel := context.WithCancel(context.Background())

	handlers, closer := config.baseHandlers()
	if closer != nil {
		defer closer.Close()
	}
	handlers.BookDownloaded = func(path string) {
		fmt.Printf("%sReceived file response.\n", clearLine)
		fmt.Println("File location: " + path)
		cancel()
	}
	handlers.BadServer = func() {
		fmt.Println("That server is not available. Try again...")
		cancel()
	}

	fmt.Printf("Sending download request.")
	config.client.StartReader(ctx, handlers)
	config.client.Download(download)
	fmt.Printf("%sSent download request.", clearLine)
	fmt.Printf("Waiting for file response.")

	registerShutdown(config.client, cancel)
	<-ctx.Done()
}

func StartSearch(config Config, query string) {
	nextSearchTime := getLastSearchTime().Add(15 * time.Second)
	instantiate(&config)
	defer config.client.Disconnect()
	ctx, cancel := context.WithCancel(context.Background())

	handlers, closer := config.baseHandlers()
	if closer != nil {
		defer closer.Close()
	}
	handlers.SearchResults = func(_ []core.BookDetail, _ []core.ParseError, path string) {
		fmt.Printf("%sReceived file response.\n", clearLine)
		fmt.Println("Results location: " + path)
		cancel()
	}
	handlers.NoResults = func() {
		fmt.Println("No results returned for your search...")
		cancel()
	}
	handlers.MatchesFound = func(num string) { fmt.Printf("Found %s search results.", num) }

	fmt.Printf("Sending search request.")
	warnIfServerOffline(query)
	time.Sleep(time.Until(nextSearchTime))

	config.client.StartReader(ctx, handlers)
	config.client.Search(config.SearchBot, query)

	setLastSearchTime()
	fmt.Printf("%sSent search request.", clearLine)
	fmt.Printf("Waiting for file response.")

	registerShutdown(config.client, cancel)
	<-ctx.Done()
}
