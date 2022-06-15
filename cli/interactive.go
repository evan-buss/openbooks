package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/evan-buss/openbooks/core"
)

func terminalMenu(config Config) {
	fmt.Print("\ns)search\ng)et book\nse)rvers\nd)one\n~> ")

	// Trim user input so we don't send 2 messages
	clean := func(message string) string { return strings.Trim(message, "\r\n") }

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = clean(input)

	switch input {
	case "s":
		fmt.Print("@search ")
		query, _ := reader.ReadString('\n')
		fmt.Println("\nSent search request.")

		nextSearchTime := getLastSearchTime().Add(15 * time.Second)
		time.Sleep(time.Until(nextSearchTime))

		core.SearchBook(config.irc, config.SearchBot, clean(query))
		setLastSearchTime()
	case "g":
		fmt.Print("Download String: ")
		message, _ := reader.ReadString('\n')
		core.DownloadBook(config.irc, clean(message))
		fmt.Println("\nSent download request.")
		warnIfServerOffline(clean(message))
	case "se":
		fmt.Println("\nAvailable Servers:")
		for _, server := range servers {
			fmt.Printf("  %s\n", server)
		}
		terminalMenu(config)
	case "d":
		fmt.Println("Disconnecting.")
		config.irc.Disconnect()
		os.Exit(0)
	default:
		fmt.Println("Invalid Selection.")
		terminalMenu(config)
	}
}

func fullHandler(config Config) core.EventHandler {
	handler := core.EventHandler{}
	addEssentialHandlers(handler, &config)

	handler[core.BadServer] = func(text string) {
		config.badServerHandler(text)
		terminalMenu(config)
	}
	handler[core.BookResult] = func(text string) {
		config.downloadHandler(text)
		terminalMenu(config)
	}
	handler[core.SearchResult] = func(text string) {
		config.searchHandler(text)
		terminalMenu(config)
	}
	handler[core.SearchAccepted] = config.searchAcceptedHandler
	handler[core.NoResults] = func(text string) {
		config.noResultsHandler(text)
		terminalMenu(config)
	}
	handler[core.MatchesFound] = config.matchesFoundHandler

	return handler
}
