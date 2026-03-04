package cmd

import (
	"os"

	"github.com/bitte-ein-bit/songbeamer-helper/log"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

const commandName = "songbeamer-helper"

// RootCmd to execute
var RootCmd = &cobra.Command{Use: commandName}

func init() {
	cobra.OnInitialize(initConfig)
}

// Execute executes the main app
func Execute() {
	RootCmd.Execute()
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			log.Fatalf("Cannot parse home dir: %s", err)
		}

		viper.SetConfigName(commandName) // name of config file (without extension)
		viper.SetConfigType("yaml")      // Set to yaml format
		viper.AddConfigPath(home)        // path to look for the config file in
		viper.AddConfigPath(".")         // optionally look for config in the working directory
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Infof("Can't read config: %v", err)
		os.Exit(1)
	}
}
