package cmd

import (
	"github.com/spf13/cobra"
)

// 登录
var syncCmd = &cobra.Command{
	Use:     "sync",
	Short:   "sync data from portainer instance to local",
	PreRunE: initAuthorizedManager,
	Long:    ``,
	RunE:    syncData,
}

func syncData(cmd *cobra.Command, args []string) error {
	// 需要初始化 token
	err := manager.SyncData()
	if err != nil {
		cmd.PrintErr("sync data", err)
		return err
	}

	cmd.Printf("sync data success\n")
	return err
}
