package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/evan-buss/openbooks/desktop"
	"github.com/evan-buss/openbooks/server"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var serverConfig server.Config

func init() {
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(desktopCmd)

	sharedFlags := pflag.NewFlagSet("", pflag.ExitOnError)
	sharedFlags.StringVarP(&serverConfig.Port, "port", "p", "5228", "Set the local network port for browser mode.")
	sharedFlags.IntP("rate-limit", "r", 10, "The number of seconds to wait between searches to reduce strain on IRC search servers. Minimum is 10 seconds.")

	serverCmd.Flags().AddFlagSet(sharedFlags)
	serverCmd.Flags().BoolVar(&serverConfig.DisableBrowserDownloads, "no-browser-downloads", false, "The browser won't recieve and download eBook files, but they are still saved to the defined 'dir' path.")
	serverCmd.Flags().StringVar(&serverConfig.Basepath, "basepath", "/", `Base path where the application is accessible. For example "/openbooks/".`)
	serverCmd.Flags().BoolVarP(&serverConfig.OpenBrowser, "browser", "b", false, "Open the browser on server start.")
	serverCmd.Flags().BoolVar(&serverConfig.Persist, "persist", false, "Persist eBooks in 'dir'. Default is to delete after sending.")
	serverCmd.Flags().StringVarP(&serverConfig.DownloadDir, "dir", "d", filepath.Join(os.TempDir(), "openbooks"), "The directory where eBooks are saved when persist enabled.")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("unable to determine HOME directory %w", err))
	}
	downloadDir := filepath.Join(homeDir, "Downloads")

	desktopCmd.Flags().AddFlagSet(sharedFlags)
	desktopCmd.Flags().StringVarP(&serverConfig.DownloadDir, "dir", "d", downloadDir, "The directory where eBooks are saved.")
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run OpenBooks in server mode.",
	Long:  "Run OpenBooks in server mode. This allows you to use a web interface to search and download eBooks.",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindGlobalFlags()
		// If cli flag isn't set (default value) check for the presence of an
		// environment variable and use it if found.
		if serverConfig.Basepath == cmd.Flag("basepath").DefValue {
			if envPath, present := os.LookupEnv("BASE_PATH"); present {
				serverConfig.Basepath = envPath
			}
		}
		ensureValidRate(cmd)
		serverConfig.Basepath = sanitizePath(serverConfig.Basepath)
	},
	Run: func(cmd *cobra.Command, args []string) {
		server.Start(serverConfig)
	},
}

var desktopCmd = &cobra.Command{
	Use:   "desktop",
	Short: "Run OpenBooks in desktop mode.",
	Long:  "Run OpenBooks in desktop mode. This allows you to run OpenBooks like a regular desktop application. This functionality utilizes your OS's native browser renderer and as such may not work on certain operating systems.",
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
}

func bindGlobalFlags() {
	serverConfig.Version = fmt.Sprintf("OpenBooks Server %s", ircVersion)
	serverConfig.UserName = globalFlags.UserName
	serverConfig.Log = globalFlags.Log
	serverConfig.Server = globalFlags.Server
	serverConfig.SearchBot = globalFlags.SearchBot
}

func ensureValidRate(cmd *cobra.Command) {
	rateLimit, _ := cmd.Flags().GetInt("rate-limit")

	// If user enters a limit that's too low, set to default of 10 seconds.
	if rateLimit < 10 {
		rateLimit = 10
	}

	serverConfig.SearchTimeout = time.Duration(rateLimit) * time.Second
}

func sanitizePath(basepath string) string {
	cleaned := path.Clean(basepath)
	if cleaned == "/" {
		return cleaned
	}
	return cleaned + "/"
}
