package cmd

import (
	"log"
	"os"

	// . "vocab/config/logging"

	"github.com/spf13/cobra"
)

var (
	confDir string = "config/config.cfg"
)

var rootCmd = &cobra.Command{
	Use:   "vocab",
	Short: "Provide services for training vocabulary",
	Long:  `Provide services for training vocabulary`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
