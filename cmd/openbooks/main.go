package main

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/evan-buss/openbooks/desktop"
	"github.com/evan-buss/openbooks/server"
	"github.com/spf13/cobra"
)

// version will always match the GitHub release versions.
var version = "4.4.0"

// We only increment ircVersion when irc admins require a fix to be made.
// They can block / permit certain version numbers. ircVersion is the current permitted
// version number.
const ircVersion = "4.3.0"

type GlobalFlags struct {
	UserName  string
	Server    string
	Log       bool
	SearchBot string
}

var globalFlags GlobalFlags

func init() {
	rootCmd.PersistentFlags().StringVarP(&globalFlags.UserName, "name", "n", generateUserName(), "Use a name that isn't randomly generated. One word only.")
	rootCmd.PersistentFlags().StringVarP(&globalFlags.Server, "server", "s", "irc.irchighway.net", "IRC server to connect to.")
	rootCmd.PersistentFlags().BoolVarP(&globalFlags.Log, "log", "l", false, "Save raw IRC logs for each client connection.")
	rootCmd.PersistentFlags().StringVar(&globalFlags.SearchBot, "searchbot", "search", "The IRC bot that handles search queries. Try 'searchook' if 'search' is down.")
}

// var rootCmd = &cobra.Command{
// 	Use:     "openbooks",
// 	Short:   "Quickly and easily download eBooks from IRCHighway.",
// 	Version: version,
// }

var rootCmd = &cobra.Command{
	Use:   "openbooks",
	Short: "Quickly and easily download eBooks from IRCHighway.",
	Long:  "Runs OpenBooks in desktop mode. This allows you to run OpenBooks like a regular desktop application. This functionality utilizes your OS's native browser renderer and as such may not work on certain operating systems.",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindGlobalFlags()
		serverConfig.DisableBrowserDownloads = true
		serverConfig.OpenBrowser = false
		serverConfig.Basepath = "/"
		serverConfig.Persist = true
	},
	Run: func(cmd *cobra.Command, args []string) {
		go server.Start(serverConfig)
		desktop.StartWebView(fmt.Sprintf("http://127.0.0.1:%s", path.Join(serverConfig.Port+serverConfig.Basepath)), false)
	},
	Version: version,
}

func main() {
	// Don't block if launched from explorer.
	cobra.MousetrapHelpText = ""

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Generate a random username to avoid IRC name collisions if multiple users are hosting
// at the same time.
func generateUserName() string {
	rand.Seed(time.Now().UnixNano())
	gofakeit.Seed(int64(rand.Int()))
	return fmt.Sprintf("%s_%s", gofakeit.Adjective(), gofakeit.Noun())
}
