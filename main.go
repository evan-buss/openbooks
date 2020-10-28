package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"strings"
	"syscall"
	"path/filepath"

	"github.com/evan-buss/openbooks/cli"
	"github.com/evan-buss/openbooks/irc"
	"github.com/evan-buss/openbooks/server"
)

var logIRC bool
var cliMode bool
var userName string
var port string
var programDir string

// Retrieve command line arguments and set appropriate variables
func init() {
	user, _ := user.Current()
	userName = strings.Split(user.Name, " ")[0]
	flag.BoolVar(&logIRC, "log", false, "Save IRC logs to irc_log.txt")
	flag.BoolVar(&cliMode, "cli", false, "Launch OpenBooks in the terminal instead of the web UI")
	flag.StringVar(&userName, "name", userName, "Use a name that differs from your account name. One word only")
	flag.StringVar(&port, "port", "5228", "Set the local network port for browser mode")
}

// Determine what mode to run the application in (CLI or Web Server)
func main() {
	flag.Parse()

	// Username can be supplied via ARGS or found from the user's system name
	if strings.Contains(userName, " ") {
		// If there is a space split it and take the first word
		fmt.Println("Using first word of entered username")
		userName = strings.Split(userName, " ")[0]
	}

	conn := irc.New(userName, userName)
	conn.Logging = logIRC

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		conn.Disconnect()
		os.Exit(1)
	}()

	if cliMode {
		cli.Start(conn)
	} else {
		programDir := filepath.Join(getProgramDir(), "downloadedEbooks")
		if _, err := os.Stat("downloadedEbooks"); os.IsNotExist(err) {
			err := os.RemoveAll(programDir)
			//fmt.Println("Directory already exists, EXTERMINATEEEE")
			if err != nil {
				panic(err)
			}
	}
	err := os.Mkdir(programDir,0777)
	if err != nil {
		panic(err)
	}
		fmt.Println("To access this go to: 127.0.0.1:"+port)
		server.Start(conn, port)
	}
}

func getProgramDir() string{
	ex, err := os.Executable()
    if err != nil {
        panic(err)
	}
	//fmt.Println("The program dir is: "+filepath.Dir(ex))
    return(filepath.Dir(ex))
}
