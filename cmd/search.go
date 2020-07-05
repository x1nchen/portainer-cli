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

// STEP
// 1. fuzzy match the container name or image name?
// 2. search from cache(boltdb) with step 1
// 3. verify the container from cache by call docker api
// 4. show container name and node name with list formation
func search(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	//
	_, err := manager.pclient.ListEndpoint(ctx)
	if err != nil {
		cmd.PrintErr(err)
		return err
	}


	// in
	res, err := manager.pclient.ListContainer(ctx, 78)
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
