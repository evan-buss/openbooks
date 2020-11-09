package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/evan-buss/openbooks/cli"
	"github.com/evan-buss/openbooks/core"
	"github.com/evan-buss/openbooks/server"
)

var config core.Config

// Determine what mode to run the application in (CLI or Web Server)
func main() {
	parseArgs()
	flag.Parse()

	// Username can be supplied via ARGS or found from the user's system name
	if strings.Contains(config.UserName, " ") {
		// If there is a space split it and take the first word
		fmt.Println("Using first word of entered username")
		config.UserName = strings.Split(config.UserName, " ")[0]
	}

	if _, isDocker := os.LookupEnv("IS_DOCKER"); isDocker {
		// Download directory is exported as a volume from the container image
		config.Log = false
		config.CliMode = false
		config.OpenBrowser = false
		config.Port = "80"
		config.DownloadDir = "/books"

		if config.UserName == "docker" {
			config.UserName = generateUserName()
			fmt.Printf("Setting IRC Name: %s\n", config.UserName)
		}
	}

	if config.CliMode {
		cli.Start(config)
	} else {
		server.Start(config)
	}
}

// Retrieve command line arguments and set appropriate variables
func parseArgs() {
	user, _ := user.Current()
	config.UserName = strings.Split(user.Name, " ")[0]

	flag.BoolVar(&config.Log, "log", false, "Save IRC logs to irc_log.txt.")
	flag.BoolVar(&config.CliMode, "cli", false, "Launch OpenBooks in the terminal instead of the web UI.")
	flag.BoolVar(&config.OpenBrowser, "browser", false, "Open the browser on server start.")
	flag.StringVar(&config.UserName, "name", config.UserName, "Use a name that differs from your account name. One word only.")
	flag.StringVar(&config.Port, "port", "5228", "Set the local network port for browser mode.")
	flag.StringVar(&config.DownloadDir, "dir", os.TempDir(), "Only applies to server mode. CLI mode downloads to the current directory.")
	flag.BoolVar(&config.Persist, "persist", false, "Keep book files in the download dir. Default is to delete after sending. (Server mode only)")
}

func generateUserName() string {
	gofakeit.Seed(0)
	return fmt.Sprintf("%s-%s", gofakeit.Adjective(), gofakeit.Noun())
}
