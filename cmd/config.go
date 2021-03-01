package cmd

import "github.com/spf13/cobra"

// TODO get/set configuration
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config to get/set configuration",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}
