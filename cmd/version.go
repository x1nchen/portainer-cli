package cmd

import (
	"github.com/spf13/cobra"
)

var version = "v0.0.4"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of portainer",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(version)
	},
}
