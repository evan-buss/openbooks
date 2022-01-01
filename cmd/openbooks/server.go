package main

import (
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/evan-buss/openbooks/server"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().BoolP("log", "l", false, "Save raw IRC logs for each client connection.")
	serverCmd.Flags().BoolP("browser", "b", false, "Open the browser on server start.")
	serverCmd.Flags().StringP("name", "n", generateUserName(), "Use a name that isn't randomly generated. One word only.")
	serverCmd.Flags().StringP("port", "p", "5228", "Set the local network port for browser mode.")
	serverCmd.Flags().StringP("dir", "d", filepath.Join(os.TempDir(), "openbooks"), "The directory where eBooks are saved when persist enabled.")
	serverCmd.Flags().Bool("persist", false, "Persist eBooks in 'dir'. Default is to delete after sending.")
	serverCmd.Flags().String("basepath", "/", `Base path where the application is accessible. For example "/openbooks/".`)
	serverCmd.Flags().StringP("server", "s", "irc.irchighway.net", "IRC server to connect to.")
	serverCmd.Flags().IntP("rate-limit", "r", 10, "The number of seconds to wait between searches to reduce strain on IRC search servers. Minimum is 10 seconds.")
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run OpenBooks in server mode.",
	Long:  "Run OpenBooks in server mode. This allows you to use a web interface to search and download eBooks.",
	Run: func(cmd *cobra.Command, args []string) {
		log, _ := cmd.Flags().GetBool("log")
		browser, _ := cmd.Flags().GetBool("browser")
		name, _ := cmd.Flags().GetString("name")
		port, _ := cmd.Flags().GetString("port")
		dir, _ := cmd.Flags().GetString("dir")
		persist, _ := cmd.Flags().GetBool("persist")
		basepath, _ := cmd.Flags().GetString("basepath")
		url, _ := cmd.Flags().GetString("server")
		rateLimit, _ := cmd.Flags().GetInt("rate-limit")

		// If cli flag isn't set (default value) check for the presence of an
		// environment variable and use it if found.
		if basepath == cmd.Flag("basepath").DefValue {
			if envPath, present := os.LookupEnv("BASE_PATH"); present {
				basepath = envPath
			}
		}

		// If user enters a limit that's too low, set to default of 10 seconds.
		if rateLimit < 10 {
			rateLimit = 10
		}

		config := server.Config{
			Log:           log,
			OpenBrowser:   browser,
			UserName:      name,
			Port:          port,
			DownloadDir:   dir,
			Persist:       persist,
			Basepath:      sanitizePath(basepath),
			Server:        url,
			SearchTimeout: time.Duration(rateLimit) * time.Second,
		}

		server.Start(config)
	},
}

func sanitizePath(basepath string) string {
	cleaned := path.Clean(basepath)
	if cleaned == "/" {
		return cleaned
	}
	return cleaned + "/"
}
