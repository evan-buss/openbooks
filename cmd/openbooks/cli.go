package main

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/evan-buss/openbooks/cli"
	"github.com/spf13/cobra"
)

func newCliCommand(globalFlags GlobalFlags) *cobra.Command {
	var cliConfig cli.Config

	cliCmd := &cobra.Command{
		Use:   "cli",
		Short: "Run openbooks from the terminal in interactive CLI mode.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cliConfig.Version = globalFlags.UserAgent
			cliConfig.UserName = globalFlags.UserName
			cliConfig.Server = globalFlags.Server
			cliConfig.Log = globalFlags.Log
			cliConfig.SearchBot = globalFlags.SearchBot
			cliConfig.EnableTLS = globalFlags.EnableTLS

			if globalFlags.Debug {
				spew.Dump(cliConfig)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			cli.StartInteractive(cliConfig)
		},
	}

	cliCmd.AddCommand(newDownloadCmd(cliConfig))
	cliCmd.AddCommand(newSearchCmd(cliConfig))

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Could not get current working directory.", err)
	}

	cliCmd.PersistentFlags().StringVarP(&cliConfig.Dir, "dir", "d", cwd, "Directory where files are downloaded.")

	return cliCmd
}

func newDownloadCmd(config cli.Config) *cobra.Command {
	downloadCmd := &cobra.Command{
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

	return downloadCmd
}

func newSearchCmd(config cli.Config) *cobra.Command {
	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "Searches for a book and exits.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cli.StartSearch(config, args[0])
		},
	}

	return searchCmd
}
