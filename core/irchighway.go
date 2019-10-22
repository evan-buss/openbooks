package core

import (
	"time"

	"github.com/evan-buss/openbooks/irc"
)

// Wrapper around the IRC connection with irc.irchighway.net specific
// BOT commands

// Join connects to the irc.irchighway.net server and joins the #ebooks channel
func Join(irc *irc.Conn) {
	irc.Connect("irc.irchighway.net")
	// Wait before joining the ebooks room
	// Often you recieve a private message from the server
	time.Sleep(time.Second * 2)
	irc.JoinChannel("ebooks")
}

// SearchBook sends a search query to the search bot
func SearchBook(irc *irc.Conn, query string) {
	irc.SendMessage("@search " + query)
}

// DownloadBook sends the book string to the download bot
func DownloadBook(irc *irc.Conn, book string) {
	irc.SendMessage(book)
}

func GetUsers(irc *irc.Conn) {
	irc.GetUsers("ebooks")
}
