package main

import (
	"flag"
	"os"
	"os/signal"
	"os/user"
	"strings"
	"syscall"

	"github.com/evan-buss/openbooks/cli"
	"github.com/evan-buss/openbooks/irc"
	"github.com/evan-buss/openbooks/server"
)

var logIRC bool
var cliMode bool
var userName string
var webPort int

// Retrieve command line arguments and set appropriate variables
func init() {
	user, _ := user.Current()
	userName = user.Name
	flag.BoolVar(&logIRC, "log", false, "Save IRC logs")
	flag.BoolVar(&cliMode, "cli", false, "Launch OpenBooks in the terminal instead of the web UI")
	flag.StringVar(&userName, "name", userName, "Use a name that differs from your account name.")
	flag.IntVar(&webPort, "port", 8080, "Use a custom port for website server")
}

// Determine what mode to run the application in (CLI or Web Server)
func main() {
	flag.Parse()

	// Username can be supplied via ARGS or found from the user's system name
	if strings.Contains(userName, " ") {
		// If there is a space split it and take the first word
		userName = strings.Split(userName, " ")[0]
	}

	irc := irc.New(userName, userName)
	irc.Logging = logIRC

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		irc.Disconnect()
		os.Exit(1)
	}()

	if cliMode {
		cli.Start(irc)
	} else {
		server.Start(irc, webPort)
	}
}
