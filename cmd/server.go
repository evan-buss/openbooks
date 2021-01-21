package cmd

import (
	"fmt"
	"os"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/evan-buss/openbooks/server"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.LocalFlags().BoolP("log", "l", false, "Save IRC logs to irc_log.txt.")
	serverCmd.LocalFlags().BoolP("browser", "b", false, "Open the browser on server start.")
	serverCmd.LocalFlags().StringP("name", "n", userName, "Use a name that differs from your account name. One word only.")
	serverCmd.LocalFlags().String("port", "5228", "Set the local network port for browser mode.")
	serverCmd.LocalFlags().StringP("dir", "d", os.TempDir(), "The directory where files are saved when persist enabled.")
	serverCmd.LocalFlags().BoolP("persist", "p", false, "Persist eBooks in 'dir'. Default is to delete after sending.")
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run OpenBooks in server mode.",
	Run: func(cmd *cobra.Command, args []string) {
		log, _ := cmd.LocalFlags().GetBool("log")
		browser, _ := cmd.LocalFlags().GetBool("browser")
		name, _ := cmd.LocalFlags().GetString("name")
		port, _ := cmd.LocalFlags().GetString("port")
		dir, _ := cmd.LocalFlags().GetString("dir")
		persist, _ := cmd.LocalFlags().GetBool("persist")

		config := server.Config{
			Log:         log,
			OpenBrowser: browser,
			UserName:    name,
			Port:        port,
			DownloadDir: dir,
			Persist:     persist,
		}

		server.Start(config)
	},
}

// if _, isDocker := os.LookupEnv("IS_DOCKER"); isDocker {
// 	// Download directory is exported as a volume from the container image
// 	config.Log = false
// 	config.CliMode = false
// 	config.OpenBrowser = false
// 	config.Port = "80"
// 	config.DownloadDir = "/books"

// 	if config.UserName == "docker" {
// 		config.UserName = generateUserName()
// 		fmt.Printf("Setting IRC Name: %s\n", config.UserName)
// 	}
// }

func generateUserName() string {
	gofakeit.Seed(0)
	return fmt.Sprintf("%s-%s", gofakeit.Adjective(), gofakeit.Noun())
}
