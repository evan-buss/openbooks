package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/evan-buss/openbooks/irc"
)

// Specific irc.irchighway.net commands
// Also works with Unanet as they use the same operator software

// Join connects to the irc.irchighway.net server and joins the #ebooks channel
func Join(irc *irc.Conn, address string, channel string, enableTLS bool) error {
	err := irc.Connect(address, enableTLS)
	if err != nil {
		return err
	}
	// Wait before joining the ebooks room
	// Often you recieve a private message from the server
	time.Sleep(time.Second * 2)
	irc.JoinChannel(channel)
	return nil
}

// SearchBook sends a search query to the search bot
func SearchBook(irc *irc.Conn, searchBot string, query string) {
	searchBot = strings.TrimPrefix(searchBot, "@")
	irc.SendMessage(fmt.Sprintf("@%s %s", searchBot, query))
}

// DownloadBook sends the book string to the download bot
func DownloadBook(irc *irc.Conn, book string) {
	irc.SendMessage(book)
}

// SendVersionInfo sends a CTCP Version response
func SendVersionInfo(irc *irc.Conn, line string, version string) {
	// Line format is like ":messager PRIVMSG #channel: message"
	// we just want the messager without the colon
	sender := strings.Split(line, " ")[0][1:]
	// TODO: Figure out if there's an automated way to adjust this...
	irc.SendNotice(sender, fmt.Sprintf("\x01%s\x01", version))
}
