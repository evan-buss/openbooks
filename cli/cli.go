package cli

import (
	"context"
	"fmt"

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

	conn := instantiate(config)
	config.irc = conn

	ctx, cancel := context.WithCancel(context.Background())
	registerShutdown(conn, cancel)

	handler := fullHandler(config)
	if config.Log {
		file := config.setupLogger(handler)
		defer file.Close()
	}

	go core.StartReader(ctx, conn, handler)
	terminalMenu(conn)

	<-ctx.Done()
}

func StartDownload(config Config, download string) {
	conn := instantiate(config)
	defer conn.Close()
	ctx, cancel := context.WithCancel(context.Background())

	handler := core.EventHandler{}
	handler[core.BookResult] = func(text string) {
		config.downloadHandler(text)
		cancel()
	}
	if config.Log {
		file := config.setupLogger(handler)
		defer file.Close()
	}

	go core.StartReader(ctx, conn, handler)
	fmt.Printf("Sending download request... ")
	core.DownloadBook(conn, download)
	fmt.Println("✅")

	registerShutdown(conn, cancel)
	<-ctx.Done()
}

func StartSearch(config Config, query string) {
	conn := instantiate(config)
	defer conn.Close()
	ctx, cancel := context.WithCancel(context.Background())

	handler := core.EventHandler{}
	handler[core.SearchResult] = func(text string) {
		config.searchHandler(text)
		cancel()
	}
	handler[core.MatchesFound] = config.matchesFoundHandler
	if config.Log {
		file := config.setupLogger(handler)
		defer file.Close()
	}

	go core.StartReader(ctx, conn, handler)
	fmt.Printf("Sending search request... ")
	core.SearchBook(conn, query)
	fmt.Println("✅")

	registerShutdown(conn, cancel)
	<-ctx.Done()
}
