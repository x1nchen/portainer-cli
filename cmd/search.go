package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

func init() {
	// searchCmd.Flags().StringVarP(&User, "user", "u", "", "user")
	// searchCmd.Flags().StringVarP(&Password, "password", "p", "", "user")
}

// 登录
var searchCmd = &cobra.Command{
	Use:     "search",
	Short:   "search container with fuzzy name",
	PreRunE: initAuthorizedManager,
	Long:    ``,
	RunE:    search,
}

func search(cmd *cobra.Command, args []string) error {
	res, err := manager.pclient.ListContainer(context.Background(), 78)
	if err != nil {
		cmd.PrintErr(err)
		return err
	}
	for _, r := range res {
		// strip the heading "/"
		cmd.Println(r.Names[0][1:], r.Image, r.Status)
	}
	return err
}
