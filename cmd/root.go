package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/toudi/kwity/config"
)

var rootCmd = &cobra.Command{
	Use:   "kwity",
	Short: "KISS invoice issuer",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	var err error

	viper.SetConfigFile("config.yaml")

	if err = viper.ReadInConfig(); err != nil {
		fmt.Printf("found error in config: %v\n", err)
		os.Exit(1)
	}

	if err = viper.Unmarshal(&config.Config); err != nil {
		fmt.Printf("could not unmarshall config: %v\n", err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
