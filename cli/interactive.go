package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/evan-buss/openbooks/core"
	"github.com/evan-buss/openbooks/irc"
)

var lastSearch time.Time

func terminalMenu(irc *irc.Conn) {
	fmt.Print("\ns)search\ng)et book\nd)one\n~> ")

	// Trim user input so we don't send 2 messages
	clean := func(message string) string { return strings.Trim(message, "\r\n") }

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = clean(input)

	switch input {
	case "s":
		fmt.Print("@search ")
		message, _ := reader.ReadString('\n')
		fmt.Println("\nSent search request.")
		time.Sleep(time.Until(lastSearch.Add(time.Second * 15)))
		core.SearchBook(irc, clean(message))
		lastSearch = time.Now()
	case "g":
		fmt.Print("Download String: ")
		message, _ := reader.ReadString('\n')
		core.DownloadBook(irc, clean(message))
		fmt.Println("\nSent download request.")
	case "d":
		fmt.Println("Disconnecting.")
		irc.Disconnect()
		os.Exit(0)
	default:
		fmt.Println("Invalid Selection.")
		terminalMenu(irc)
	}
}

func fullHandler(config Config) core.EventHandler {
	handler := core.EventHandler{}
	addRequiredHandlers(handler, &config)

	handler[core.BadServer] = func(text string) {
		config.badServerHandler(text)
		terminalMenu(config.irc)
	}
	handler[core.BookResult] = func(text string) {
		config.downloadHandler(text)
		terminalMenu(config.irc)
	}
	handler[core.SearchResult] = func(text string) {
		config.searchHandler(text)
		terminalMenu(config.irc)
	}
	handler[core.SearchAccepted] = config.searchAcceptedHandler
	handler[core.NoResults] = func(text string) {
		config.noResultsHandler(text)
		terminalMenu(config.irc)
	}
	handler[core.MatchesFound] = config.matchesFoundHandler

	return handler
}
