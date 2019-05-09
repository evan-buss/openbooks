package main

import (
	"bufio"
	"fmt"
	"openbooks/dcc"
	"openbooks/irc"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Println("Welcome to OpenBooks")
	irc := irc.New("evan_bot", "Evan C Buss")
	// irc.Connect("localhost")
	irc.Connect("irc.irchighway.net")

	time.Sleep(time.Second * 2)

	irc.JoinChannel("ebooks")

	go readDaemon(irc)

	fmt.Println("Connection established...")
	fmt.Println("Server rules mandate a 30 second wait period")

	for i := 30; i > 0; i-- {
		fmt.Print("\r" + strconv.Itoa(i) + "...")
		time.Sleep(time.Second)
	}

	fmt.Print("\r")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\ns)end\nd)isconnect\n-> ")
		input, _ := reader.ReadString('\n')

		switch input {
		case "s\n":
			message, _ := reader.ReadString('\n')
			irc.SendMessage(message)
		case "d\n":
			fmt.Println("disonnecting")
			irc.Disconnect()
			os.Exit(0)
		default:
			fmt.Println("Invalid Selection")
		}
	}
}

func readDaemon(irc *irc.Conn) {
	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	f.WriteString("\n==================== NEW LOG ======================\n")

	if err != nil {
		panic(err)
	}
	defer f.Close()

	for {
		text := irc.GetMessage()
		// fmt.Println(text)
		f.WriteString(text)

		if strings.Contains(text, "DCC SEND") {
			go dcc.NewDownload(text, false)
		}
	}
}
