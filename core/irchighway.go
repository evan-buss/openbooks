package core

import (
	"time"

	"github.com/evan-buss/openbooks/irc"
)

// Specific irc.irchighway.net commands

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
	//Send the request after a delay to ensure we are actuall in the channel
	time.Sleep(time.Second*5)
	irc.GetUsers("ebooks")
}
