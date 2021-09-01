package cli

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/evan-buss/openbooks/core"
	"github.com/evan-buss/openbooks/irc"
)

// Config is used to configure CLI mode settings.
type Config struct {
	UserName string // Username to use when connecting to IRC
	Log      bool   // True if IRC messages should be logged
}

// Reader is a way to recieve input from the user
var reader *bufio.Reader

// IRC is the current IRC connection
var conn *irc.Conn

// Start instantiates the OpenBooks CLI interface
func Start(config Config) {
	conn := irc.New(config.UserName, "OpenBooks CLI")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		conn.Disconnect()
		os.Exit(1)
	}()

	fmt.Println("=======================================")
	fmt.Println("          Welcome to OpenBooks         ")
	fmt.Println("=======================================")

	core.Join(conn)

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Could not get current working directory.", err)
	}

	exitSignal := make(chan struct{})
	go core.ReadDaemon(conn, Handler{irc: conn, downloadDir: cwd}, config.Log, exitSignal)

	fmt.Println("Connection established...")

	fmt.Print("\r")

	reader = bufio.NewReader(os.Stdin)

	// Get the first input
	// Reader, IRC
	menu()
	// We make a channel to block forever. We want the reader daemon to run forever
	<-exitSignal
}

//reader *bufio.Reader, irc *irc.Conn
func menu() {
	fmt.Print("\ns)search\ng)et book\nd)one\n~> ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimRight(input, "\n")
	input = strings.TrimRight(input, "\r")

	switch input {
	case "s":
		fmt.Print("@search ")
		message, _ := reader.ReadString('\n')
		core.SearchBook(conn, message)
	case "g":
		fmt.Print("Download String: ")
		message, _ := reader.ReadString('\n')
		core.DownloadBook(conn, message)
	case "d":
		fmt.Println("Disconnecting.")
		conn.Disconnect()
		os.Exit(0)
	default:
		fmt.Println("Invalid Selection.")
		menu()
	}
}
