package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/evan-buss/openbooks/desktop"
	"github.com/evan-buss/openbooks/server"
	"github.com/spf13/cobra"
)

// version will always match the GitHub release versions.
var version = "4.4.0"

// We only increment ircVersion when irc admins require a fix to be made.
// They can block / permit certain version numbers. ircVersion is the current permitted
// version number.
var ircVersion = "4.3.0"

type GlobalFlags struct {
	UserName  string
	Server    string
	Log       bool
	SearchBot string
	EnableTLS bool
}

var debug bool
var globalFlags GlobalFlags
var desktopConfig server.Config

func init() {
	desktopCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug mode.")
	desktopCmd.PersistentFlags().StringVarP(&globalFlags.UserName, "name", "n", generateUserName(), "Use a name that isn't randomly generated. One word only.")
	desktopCmd.PersistentFlags().StringVarP(&globalFlags.Server, "server", "s", "irc.irchighway.net:6697", "IRC server to connect to.")
	desktopCmd.PersistentFlags().BoolVar(&globalFlags.EnableTLS, "tls", true, "Connect to server using TLS.")
	desktopCmd.PersistentFlags().BoolVarP(&globalFlags.Log, "log", "l", false, "Save raw IRC logs for each client connection.")
	desktopCmd.PersistentFlags().StringVar(&globalFlags.SearchBot, "searchbot", "search", "The IRC bot that handles search queries. Try 'searchook' if 'search' is down.")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("unable to determine $HOME directory %w", err))
	}
	downloadDir := filepath.Join(homeDir, "Downloads")

	desktopCmd.Flags().StringVarP(&desktopConfig.Port, "port", "p", "5228", "Set the local network port for browser mode.")
	desktopCmd.Flags().IntP("rate-limit", "r", 10, "The number of seconds to wait between searches to reduce strain on IRC search servers. Minimum is 10 seconds.")
	desktopCmd.Flags().StringVarP(&desktopConfig.DownloadDir, "dir", "d", downloadDir, "The directory where eBooks are saved.")
}

var desktopCmd = &cobra.Command{
	Use:   "openbooks",
	Short: "Quickly and easily download eBooks from IRCHighway.",
	Long:  "Runs OpenBooks in desktop mode. This allows you to run OpenBooks like a regular desktop application. This functionality utilizes your OS's native browser renderer and as such may not work on certain operating systems.",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindGlobalServerFlags(&desktopConfig)
		rateLimit, _ := cmd.Flags().GetInt("rate-limit")
		ensureValidRate(rateLimit, &desktopConfig)
		desktopConfig.DisableBrowserDownloads = true
		desktopConfig.Basepath = "/"
		desktopConfig.Persist = true
	},
	Run: func(cmd *cobra.Command, args []string) {
		if debug {
			spew.Dump(desktopConfig)
		}

		go server.Start(desktopConfig)
		desktop.StartWebView(fmt.Sprintf("http://127.0.0.1:%s", path.Join(desktopConfig.Port+desktopConfig.Basepath)), debug)
	},
	Version: version,
}

func main() {
	// Don't block if launched from explorer.
	cobra.MousetrapHelpText = ""

	if err := desktopCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
