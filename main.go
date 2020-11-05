package main

import (
	"flag"
	"fmt"
	"os/user"
	"strings"

	"github.com/evan-buss/openbooks/cli"
	"github.com/evan-buss/openbooks/core"
	"github.com/evan-buss/openbooks/server"
)

var cliMode bool
var config core.Config

// Retrieve command line arguments and set appropriate variables
func init() {
	user, _ := user.Current()
	config.UserName = strings.Split(user.Name, " ")[0]
	flag.BoolVar(&config.Log, "log", false, "Save IRC logs to irc_log.txt.")
	flag.BoolVar(&cliMode, "cli", false, "Launch OpenBooks in the terminal instead of the web UI.")
	flag.BoolVar(&config.OpenBrowser, "browser", false, "Open the browser on server start.")
	flag.StringVar(&config.UserName, "name", config.UserName, "Use a name that differs from your account name. One word only.")
	flag.StringVar(&config.Port, "port", "5228", "Set the local network port for browser mode.")
}

// Determine what mode to run the application in (CLI or Web Server)
func main() {
	flag.Parse()

	// Username can be supplied via ARGS or found from the user's system name
	if strings.Contains(config.UserName, " ") {
		// If there is a space split it and take the first word
		fmt.Println("Using first word of entered username")
		config.UserName = strings.Split(config.UserName, " ")[0]
	}

	if cliMode {
		cli.Start(config)
	} else {
		server.Start(config)
	}
}
