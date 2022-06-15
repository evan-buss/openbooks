package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/evan-buss/openbooks/server"

	"github.com/spf13/cobra"
)

var serverConfig server.Config

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().BoolVarP(&serverConfig.OpenBrowser, "browser", "b", false, "Open the browser on server start.")
	serverCmd.Flags().StringVarP(&serverConfig.Port, "port", "p", "5228", "Set the local network port for browser mode.")
	serverCmd.Flags().BoolVar(&serverConfig.Persist, "persist", false, "Persist eBooks in 'dir'. Default is to delete after sending.")
	serverCmd.Flags().StringVar(&serverConfig.Basepath, "basepath", "/", `Base path where the application is accessible. For example "/openbooks/".`)
	serverCmd.Flags().IntP("rate-limit", "r", 10, "The number of seconds to wait between searches to reduce strain on IRC search servers. Minimum is 10 seconds.")
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run OpenBooks in server mode.",
	Long:  "Run OpenBooks in server mode. This allows you to use a web interface to search and download eBooks.",
	PreRun: func(cmd *cobra.Command, args []string) {
		serverConfig.Version = fmt.Sprintf("OpenBooks Server %s", ircVersion)
		serverConfig.UserName = globalFlags.UserName
		serverConfig.DownloadDir = globalFlags.DownloadDir
		serverConfig.Log = globalFlags.Log
		serverConfig.Server = globalFlags.Server
		serverConfig.SearchBot = globalFlags.SearchBot
	},
	Run: func(cmd *cobra.Command, args []string) {
		rateLimit, _ := cmd.Flags().GetInt("rate-limit")

		// If cli flag isn't set (default value) check for the presence of an
		// environment variable and use it if found.
		if serverConfig.Basepath == cmd.Flag("basepath").DefValue {
			if envPath, present := os.LookupEnv("BASE_PATH"); present {
				serverConfig.Basepath = envPath
			}
		}

		// If user enters a limit that's too low, set to default of 10 seconds.
		if rateLimit < 10 {
			rateLimit = 10
		}

		serverConfig.SearchTimeout = time.Duration(rateLimit) * time.Second
		serverConfig.Basepath = sanitizePath(serverConfig.Basepath)

		server.Start(serverConfig)
	},
}

func sanitizePath(basepath string) string {
	cleaned := path.Clean(basepath)
	if cleaned == "/" {
		return cleaned
	}
	return cleaned + "/"
}
