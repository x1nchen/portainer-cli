package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of portainer",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("v0.0.1")
	},
}