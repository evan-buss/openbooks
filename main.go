package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/evan-buss/openbooks/dcc"
	"github.com/evan-buss/openbooks/irc"
)

var logIRC bool
var userName string

// Retrieve command line arguments and set appropriate variables
func parseArgs() {
	user, _ := user.Current()
	userName = user.Name

	flag.BoolVar(&logIRC, "log", false, "Save IRC logs")
	flag.StringVar(&userName, "name", userName, "Use a name that differs from your account name.")
	flag.Parse()
}

// Establishes connection with IRC server, displays menu, spawns goroutines
func main() {
	fmt.Println("=======================================")
	fmt.Println("          Welcome to OpenBooks         ")
	fmt.Println("=======================================")

	parseArgs()

	// Username can be supplied via ARGS or found from the user's system name
	if strings.Contains(userName, " ") {
		log.Fatal("Please supply a single word username. Cannot use " + userName)
	}

	irc := irc.New(userName, userName)
	irc.Connect("irc.irchighway.net")

	// Wait before joining the ebooks room
	// Often you recieve a private message from the server
	time.Sleep(time.Second * 2)
	irc.JoinChannel("ebooks")

	statusC := make(chan bool)
	stateC := make(chan bool)

	go readDaemon(irc, statusC, stateC)

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

// Designed to be launched as a goroutine. Listens for specific messages and
// 	responds accordingly
// Params: irc - IRC connection
//  			 statusC - boolean channel returns true when the download has finished
// 				 stateC - boolean channel recieves messages from user menu
// 									(true = book download, false = search results download)
func readDaemon(irc *irc.Conn, statusC chan<- bool, stateC <-chan bool) {
	var f *os.File
	var err error

	if logIRC {
		f, err = os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
		f.WriteString("\n==================== NEW LOG ======================\n")

		if err != nil {
			panic(err)
		}
	}

	isBook := false

	doneChan := make(chan bool)

	for {
		text := irc.GetMessage()

		if logIRC {
			f.WriteString(text)
		}

		select {
		// Get state of app from menu. isBook = download book; !isBook = search
		case isBook = <-stateC:
		default:
		}

		if strings.Contains(text, "DCC SEND") {
			// Respond to Direct Client-to-Client downloads
			go dcc.NewDownload(text, isBook, doneChan)
		} else if strings.Contains(text, "NOTICE") {
			// Notice is a message sent directly to the user
			if strings.Contains(text, "Sorry") {
				// There were no results
				fmt.Println("No results returned for that search...")
				statusC <- !isBook
			} else if strings.Contains(text, "try another server") {
				fmt.Println("That server is not available. Try again...")
				statusC <- !isBook
			} else if strings.Contains(text, "has been accepted") {
				fmt.Println("Search has been accepted. Please wait.")
			} else if strings.Contains(text, "matches") {
				start := strings.LastIndex(text, "returned") + 9
				end := strings.LastIndex(text, "matches") - 1
				fmt.Println("Your search returned " + text[start:end] + " matches.")
			}
		}
		select {
		case <-doneChan:
			// Send message when finished downloading
			statusC <- true
		default:
		}
	}
}
