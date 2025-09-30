package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var version = "v0.3.2"

var rootCmd = &cobra.Command{
	Use:     "craft",
	Short:   "Craft is a CLI tool for generating project files",
	Long:    `Craft is a command-line tool to generate services, repositories, and more for your Golang projects.`,
	Version: version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("ERROR: %v\n", err)
	}
}
