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

type GlobalFlags struct {
	UserName  string
	Server    string
	Log       bool
	SearchBot string
	EnableTLS bool
	UserAgent string
	Debug     bool
}

func newRootCommand() *cobra.Command {
	var globalFlags GlobalFlags
	var desktopConfig server.Config

	// version will always match the GitHub release versions.
	var version = "4.5.0"

	// We only increment ircVersion when irc admins require a fix to be made.
	// They can block / permit certain version numbers. ircVersion is the current permitted
	// version number.
	var ircVersion = "4.3.0"

	rootCmd := &cobra.Command{
		Use:   "openbooks",
		Short: "Quickly and easily download eBooks from IRCHighway.",
		Long:  "Runs OpenBooks in desktop mode. This allows you to run OpenBooks like a regular desktop application. This functionality utilizes your OS's native browser renderer and as such may not work on certain operating systems.",
		PreRun: func(cmd *cobra.Command, args []string) {
			bindGlobalServerFlags(&desktopConfig, globalFlags)
			rateLimit, _ := cmd.Flags().GetInt("rate-limit")
			ensureValidRate(rateLimit, &desktopConfig)
			desktopConfig.DisableBrowserDownloads = true
			desktopConfig.Basepath = "/"
			desktopConfig.Persist = true
		},
		Run: func(cmd *cobra.Command, args []string) {
			if globalFlags.Debug {
				spew.Dump(desktopConfig)
			}

			go server.Start(desktopConfig)
			desktop.StartWebView(fmt.Sprintf("http://127.0.0.1:%s", path.Join(desktopConfig.Port+desktopConfig.Basepath)), globalFlags.Debug)
		},
		Version: version,
	}

	rootCmd.PersistentFlags().StringVarP(&globalFlags.UserName, "name", "n", "", "Username used to connect to IRC server.")
	err := rootCmd.MarkPersistentFlagRequired("name")
	if err != nil {
		panic(err)
	}
	rootCmd.PersistentFlags().StringVarP(&globalFlags.Server, "server", "s", "irc.irchighway.net:6697", "IRC server to connect to.")
	rootCmd.PersistentFlags().BoolVar(&globalFlags.EnableTLS, "tls", true, "Connect to server using TLS.")
	rootCmd.PersistentFlags().BoolVarP(&globalFlags.Log, "log", "l", false, "Save raw IRC logs for each client connection.")
	rootCmd.PersistentFlags().StringVar(&globalFlags.SearchBot, "searchbot", "search", "The IRC bot that handles search queries. Try 'searchook' if 'search' is down.")
	rootCmd.PersistentFlags().StringVarP(&globalFlags.UserAgent, "useragent", "u", fmt.Sprintf("OpenBooks %s", ircVersion), "UserAgent / Version Reported to IRC Server.")
	rootCmd.PersistentFlags().BoolVar(&globalFlags.Debug, "debug", false, "Enable debug mode.")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("unable to determine $HOME directory %w", err))
	}
	downloadDir := filepath.Join(homeDir, "Downloads")

	rootCmd.Flags().StringVarP(&desktopConfig.Port, "port", "p", "5228", "Set the local network port for browser mode.")
	rootCmd.Flags().IntP("rate-limit", "r", 10, "The number of seconds to wait between searches to reduce strain on IRC search servers. Minimum is 10 seconds.")
	rootCmd.Flags().StringVarP(&desktopConfig.DownloadDir, "dir", "d", downloadDir, "The directory where eBooks are saved.")

	rootCmd.AddCommand(newCliCommand(globalFlags))
	rootCmd.AddCommand(newServerCommand(globalFlags))

	return rootCmd
}

func main() {
	// Don't block if launched from explorer.
	cobra.MousetrapHelpText = ""

	if err := newRootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
