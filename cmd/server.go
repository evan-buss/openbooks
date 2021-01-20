package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var name string

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.PersistentFlags().StringVarP(&name, "name", "n", "evan", "Set the name for the child.")
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run OpenBooks in server mode.",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
		fmt.Printf("His name was... %v", name)
	},
}
