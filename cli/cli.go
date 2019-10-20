package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/evan-buss/openbooks/core"
	"github.com/evan-buss/openbooks/irc"
)

// Start instantiates the OpenBooks CLI interface
func Start(irc *irc.Conn) {
	fmt.Println("=======================================")
	fmt.Println("          Welcome to OpenBooks         ")
	fmt.Println("=======================================")

	irc.Connect("irc.irchighway.net")
	// Wait before joining the ebooks room
	// Often you recieve a private message from the server
	time.Sleep(time.Second * 2)
	irc.JoinChannel("ebooks")

	statusC := make(chan bool)
	stateC := make(chan bool)

	go core.ReadDaemon(irc, statusC, stateC)

	fmt.Println("Connection established...")

	for i := 30; i > 0; i-- {
		fmt.Print("\rServer rules mandate a " + strconv.Itoa(i) + " second wait period   ")
		time.Sleep(time.Second)
	}

	fmt.Print("\r")

	reader := bufio.NewReader(os.Stdin)
	isDownloading := false // Hide menu when download in progress

	for {
		select {
		case <-statusC:
			isDownloading = false
		default:
		}

		if !isDownloading {
			fmt.Print("\ns)search\ng)et book\nd)one\n~> ")

			input, _ := reader.ReadString('\n')
			input = strings.TrimRight(input, "\n")
			input = strings.TrimRight(input, "\r")

			switch input {
			case "s":
				fmt.Print("@search ")
				message, _ := reader.ReadString('\n')
				irc.SendMessage("@search " + message)
				stateC <- false
				isDownloading = true
			case "g":
				fmt.Print("Download String: ")
				message, _ := reader.ReadString('\n')
				irc.SendMessage(message)
				stateC <- true
				isDownloading = true
			case "d":
				fmt.Println("disonnecting")
				irc.Disconnect()
				os.Exit(0)
			default:
				fmt.Println("Invalid Selection")
			}
		}
		time.Sleep(time.Millisecond * 100)
	}
}
