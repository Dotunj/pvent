package main

import (
	"fmt"

	"github.com/dotunj/pvent"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Pvent Version %s\n", pvent.GetVersion())
		return nil
	},
}
