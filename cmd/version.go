package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var (
	version = "0.0.1"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Pvent",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Pvent Version %s\n", version)
		return nil
	},
}
