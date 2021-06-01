package cmd

import (
	"github.com/spf13/cobra"
)

func init() {}

// upgrade
var upgradeCmd = &cobra.Command{
	Use:     "upgrade ",
	Short:   "init with terminal box to do deploy interactively",
	PreRunE: initAuthorizedManager,
	Long:    ``,
	RunE:    upgrade,
}

func upgrade(cmd *cobra.Command, args []string) error {

	return nil
}
