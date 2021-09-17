package cli

import (
	"fmt"

	"github.com/evan-buss/openbooks/core"
)

func progress(current, total int64) {
	fmt.Printf("Progress %%%.2f\r", (float64(current)/float64(total))*100)
}

// DownloadSearchResults downloads the search results
// and sends user a response message
func (c Config) searchHandler(text string) {
	fmt.Printf("Downloading... ")
	extractedPath, err := core.DownloadExtractDCCString(c.Dir, text, progress)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("✅")
	fmt.Println("Results location: " + extractedPath)
}

// DownloadBookFile downloads the search results and sends
// a user a response message
func (c Config) downloadHandler(text string) {
	fmt.Printf("Downloading... ")
	extractedPath, err := core.DownloadExtractDCCString(c.Dir, text, progress)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("✅")
	fmt.Println("File location: " + extractedPath)
}

// NoResults is called when the user searches for something that
// is not available sends a CLI message
func (c Config) noResultsHandler(_ string) {
	fmt.Println("No results returned for your search...")
}

// BadServer is called when the user tries to download a file from a
// server that is not available.
func (c Config) badServerHandler(_ string) {
	fmt.Println("That server is not available. Try again...")
}

// SearchAccepted is called when the search has been accepted but the user
// must wait in the queue for the search to be executed.
func (c Config) searchAcceptedHandler(_ string) {
	fmt.Println("Search has been accepted. Please wait.")
}

// MatchesFound is called when the search returns the number of results
// found. Server sends the client a status update
func (c Config) matchesFoundHandler(num string) {
	fmt.Println("📚 Search returned " + num + " matches.")
}

func (c Config) pingHandler(_ string) {
	c.irc.Pong(c.Server)
}

// func (h Handler) ServerList(servers core.IrcServers) {}
