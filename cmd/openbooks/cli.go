package main

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/evan-buss/openbooks/cli"
	"github.com/spf13/cobra"
)

var config cli.Config

func init() {
	rootCmd.AddCommand(cliCmd)
	cliCmd.AddCommand(downloadCmd)
	cliCmd.AddCommand(searchCmd)

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Could not get current working directory.", err)
	}

	cliCmd.PersistentFlags().StringVarP(&config.UserName, "name", "n", generateUserName(), "Use a name that isn't randomly generated. One word only.")
	cliCmd.PersistentFlags().StringVarP(&config.Dir, "directory", "d", cwd, "Directory where files are downloaded.")
	cliCmd.PersistentFlags().BoolVarP(&config.Log, "log", "l", false, "Whether or not to log IRC messages to an output file.")
	cliCmd.PersistentFlags().StringVarP(&config.Server, "server", "s", "irc.irchighway.net", "IRC server to connect to.")
}

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Run openbooks from the terminal in CLI mode.",
	Run: func(cmd *cobra.Command, args []string) {
		cli.StartInteractive(config)
	},
}

var downloadCmd = &cobra.Command{
	Use:     "download [flags] identifier",
	Short:   "Downloads a single file and exits.",
	Example: `openbooks cli download '!Oatmeal - F. Scott Fitzgerald - The Great Gatsby.epub'`,
	Args: func(cmd *cobra.Command, args []string) error {
		err := cobra.ExactArgs(1)(cmd, args)
		if err != nil {
			return err
		}
		if !strings.HasPrefix(args[0], "!") {
			return errors.New("identifier must begin with '!'")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		cli.StartDownload(config, args[0])
	},
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Searches for a book and exits.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cli.StartSearch(config, args[0])
	},
}
