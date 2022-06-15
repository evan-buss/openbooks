package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/evan-buss/openbooks/cli"
	"github.com/spf13/cobra"
)

var cliConfig cli.Config

func init() {
	rootCmd.AddCommand(cliCmd)
	cliCmd.AddCommand(downloadCmd)
	cliCmd.AddCommand(searchCmd)

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Could not get current working directory.", err)
	}

	// Change the default to CWD for CLI mode.
	cliCmd.InheritedFlags().StringVarP(&cliConfig.Dir, "dir", "d", cwd, "Directory where files are downloaded.")
	cliCmd.PersistentFlags().StringVarP(&cliConfig.Dir, "dir", "d", cwd, "Directory where files are downloaded.")
}

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Run openbooks from the terminal in interactive CLI mode.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cliConfig.Version = fmt.Sprintf("OpenBooks CLI %s", ircVersion)
		cliConfig.UserName = globalFlags.UserName
		cliConfig.Server = globalFlags.Server
		cliConfig.Log = globalFlags.Log
		cliConfig.SearchBot = globalFlags.SearchBot
	},
	Run: func(cmd *cobra.Command, args []string) {
		cli.StartInteractive(cliConfig)
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
		cli.StartDownload(cliConfig, args[0])
	},
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Searches for a book and exits.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cli.StartSearch(cliConfig, args[0])
	},
}
