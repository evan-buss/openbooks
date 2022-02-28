package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/evan-buss/openbooks/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	cmd := NewRootCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// Build the cobra command that handles our command line tool.
func NewRootCommand() *cobra.Command {
	var logOutput bool
	var config Config

	// Define our command
	rootCmd := &cobra.Command{
		Use:   "openbooks_api",
		Short: "A service that augments OpenBooks to provide more robust searching capabilities.",
		Long: `IRC admins have issues with the rate of searches. This service will cache 
			   the latest book inventory and allow users to search directly from the HTTP 
			   API instead of going through the @search IRC operator.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			v := viper.New()

			v.SetConfigName("openbooksapi")

			v.AddConfigPath(".")
			v.AddConfigPath("$HOME")
			v.AutomaticEnv()
			util.BindFlags(cmd, v, "OBA_")

			if err := v.ReadInConfig(); err != nil {
				// It's okay if there isn't a config file
				if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
					return err
				}
			}

			err := v.Unmarshal(&config)
			if err != nil {
				return err
			}

			return nil

		},
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			out := cmd.OutOrStdout()

			if logOutput {
				currentConfig, _ := json.MarshalIndent(config, "", "  ")
				fmt.Fprintf(out, "%s\n", currentConfig)
			}

			NewServer(config)
		},
	}

	// Define cobra flags, the default value has the lowest (least significant) precedence
	rootCmd.Flags().BoolVar(&logOutput, "log", false, "Log the current config values on startup.")
	rootCmd.Flags().StringVarP(&config.Port, "port", "p", "80", "Set the port to listen on.")

	return rootCmd
}
