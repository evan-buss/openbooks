package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/evan-buss/openbooks/server"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
	userName := generateUserName()

	serverCmd.Flags().BoolP("log", "l", false, "Save IRC logs to irc_log.txt.")
	serverCmd.Flags().BoolP("browser", "b", false, "Open the browser on server start.")
	serverCmd.Flags().StringP("name", "n", userName, "Use a name that isn't randomly generated. One word only.")
	serverCmd.Flags().StringP("port", "p", "5228", "Set the local network port for browser mode.")
	serverCmd.Flags().StringP("dir", "d", os.TempDir(), "The directory where eBooks are saved when persist enabled.")
	serverCmd.Flags().Bool("persist", false, "Persist eBooks in 'dir'. Default is to delete after sending.")
	serverCmd.Flags().String("basepath", "/", `Base path where the application is accessible. For example "/openbooks/".`)
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

		// If cli flag isn't set (default value) check for the presence of an
		// environment variable and use it if found.
		if basepath == cmd.Flag("basepath").DefValue {
			if envPath, present := os.LookupEnv("BASE_PATH"); present {
				basepath = envPath
			}
		}

		config := server.Config{
			Log:         log,
			OpenBrowser: browser,
			UserName:    name,
			Port:        port,
			DownloadDir: dir,
			Persist:     persist,
			Basepath:    basepath,
		}

		server.Start(config)
	},
}

// Generate a random username to avoid IRC name collisions if multiple users are hosting
// at the same time.
func generateUserName() string {
	rand.Seed(time.Now().UnixNano())
	gofakeit.Seed(int64(rand.Int()))
	return fmt.Sprintf("%s-%s", gofakeit.Adjective(), gofakeit.Noun())
}
