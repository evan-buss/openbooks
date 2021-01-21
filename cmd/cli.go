package cmd

import (
	"os/user"
	"strings"

	"github.com/evan-buss/openbooks/cli"
	"github.com/spf13/cobra"
)

func init() {
	user, _ := user.Current()
	userName := strings.Split(user.Name, " ")[0]

	rootCmd.AddCommand(cliCmd)
	cliCmd.LocalFlags().StringP("username", "u", userName, "Username to connect to the IRC server as.")
	cliCmd.LocalFlags().BoolP("log", "l", false, "Whether or not to log IRC data to an output file.")
}

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Run openbooks from the terminal in CLI mode.",
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.LocalFlags().GetString("username")
		log, _ := cmd.LocalFlags().GetBool("log")

		config := cli.Config{
			UserName: username,
			Log:      log,
		}

		cli.Start(config)
	},
}
