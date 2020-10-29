package cli

import (
	"bufio"
	"fmt"
	"os"
	//"strconv"
	"strings"
	//"time"

	"github.com/evan-buss/openbooks/core"
	"github.com/evan-buss/openbooks/irc"
)

// Reader is a way to recieve input from the user
var Reader *bufio.Reader

// IRC is the current IRC connection
var IRC *irc.Conn

// Start instantiates the OpenBooks CLI interface
func Start(irc *irc.Conn) {
	fmt.Println("=======================================")
	fmt.Println("          Welcome to OpenBooks         ")
	fmt.Println("=======================================")

	IRC = irc
	core.Join(IRC)

	exitSignal := make(chan struct{})
	go core.ReadDaemon(irc, Handler{})

	fmt.Println("Connection established...")

	/*for i := 30; i > 0; i-- {
		fmt.Print("\rServer rules mandate a " + strconv.Itoa(i) + " second wait period   ")
		time.Sleep(time.Second)
	}*/

	fmt.Print("\r")

	Reader = bufio.NewReader(os.Stdin)

	// Get the first input
	userInput(Reader, IRC)
	// We make a channel to block forever. We want the reader daemon to run forever
	<-exitSignal
}

func userInput(reader *bufio.Reader, irc *irc.Conn) {
	fmt.Print("\ns)search\ng)et book\nd)one\n~> ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimRight(input, "\n")
	input = strings.TrimRight(input, "\r")

	switch input {
	case "s":
		fmt.Print("@search ")
		message, _ := reader.ReadString('\n')
		core.SearchBook(irc, message)
	case "g":
		fmt.Print("Download String: ")
		message, _ := reader.ReadString('\n')
		core.DownloadBook(irc, message)
	case "d":
		fmt.Println("disonnecting")
		irc.Disconnect()
		os.Exit(0)
	default:
		fmt.Println("Invalid Selection")
		userInput(reader, irc)
	}
}
