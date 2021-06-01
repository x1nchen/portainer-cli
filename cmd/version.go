package cmd

import (
	"github.com/spf13/cobra"
)

var (
	Repo = ""
	Branch = ""
	Commit = ""
	Version = ""
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of portainer",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("repo: %s \nBranch: %s\nCommit: %s\nVersion: %s\n", Repo, Branch, Commit, Version)
	},
}
