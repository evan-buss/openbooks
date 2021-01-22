package cmd

import (
	"os/user"
	"strings"

	"github.com/evan-buss/openbooks/cli"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cliCmd)

	user, _ := user.Current()
	userName := strings.Split(user.Name, " ")[0]
	cliCmd.Flags().StringP("name", "n", userName, "Username to connect to the IRC server as.")
	cliCmd.Flags().BoolP("log", "l", false, "Whether or not to log IRC data to an output file.")
}

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Run openbooks from the terminal in CLI mode.",
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		log, _ := cmd.Flags().GetBool("log")

		config := cli.Config{
			UserName: username,
			Log:      log,
		}

		cli.Start(config)
	},
}
