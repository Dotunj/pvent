package main

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "pvent",
		Short: "Pvent is a CLI tool for sending events across message brokers",
	}
)

func init() {
	rootCmd.AddCommand(dispatcherCmd)
	rootCmd.AddCommand(versionCmd)
}
