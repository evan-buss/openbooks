package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/evan-buss/openbooks/irc"
)

// Specific irc.irchighway.net commands

// Join connects to the irc.irchighway.net server and joins the #ebooks channel
func Join(irc *irc.Conn, url string) error {
	err := irc.Connect(url)
	if err != nil {
		return err
	}
	// Wait before joining the ebooks room
	// Often you recieve a private message from the server
	time.Sleep(time.Second * 2)
	irc.JoinChannel("ebooks")
	return nil
}

// SearchBook sends a search query to the search bot
func SearchBook(irc *irc.Conn, query string) {
	irc.SendMessage("@search " + query)
}

// DownloadBook sends the book string to the download bot
func DownloadBook(irc *irc.Conn, book string) {
	irc.SendMessage(book)
}

// Send a CTCP Version response
func SendVersionInfo(irc *irc.Conn, line string, version string) {
	// Line format is like ":messager PRIVMSG #channel: message"
	// we just want the messager without the colon
	sender := strings.Split(line, " ")[0][1:]
	// TODO: Figure out if there's an automated way to adjust this...
	irc.SendNotice(sender, fmt.Sprintf("\x01%s\x01", version))
}
