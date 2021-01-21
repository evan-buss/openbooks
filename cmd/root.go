package cmd

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/spf13/cobra"
)

var userName string

func init() {
	user, _ := user.Current()
	userName = strings.Split(user.Name, " ")[0]
}

var rootCmd = &cobra.Command{
	Use:   "openbooks",
	Short: "Quickly and easily download eBooks from IRCHighway.",
}

// Execute starts the root command handler.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
