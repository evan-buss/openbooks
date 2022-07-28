package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var configPath string
var config = &Config{}

var rootCmd = &cobra.Command{
	Use:     "service",
	Short:   "Service defines an OpenBooks compatible search service",
	Example: "service --config config.json",
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetConfigFile(configPath)

		viper.ReadInConfig()

		viper.AutomaticEnv()

		err := viper.UnmarshalExact(config)
		if err != nil {
			log.Fatal(err)
		}

		viper.WatchConfig()
		viper.OnConfigChange(func(in fsnotify.Event) {
			log.Println("Config changed")
			viper.UnmarshalExact(config)
		})

		if config.Debug {
			spew.Dump(config)
		}

		NewServer(config)
	},
}

func main() {
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "config.json", "Config file to use.")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
