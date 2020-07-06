package cmd

import (
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
	name := args[0]
	containers, err := manager.store.ContainerService.FuzzyFindContainerByName(name)
	if err != nil {
		cmd.PrintErr(err)
		return err
	}

	for _, container := range containers {
		cmd.Println(container.Names[0][1:], container.EndpointName, container.State)
	}

	return nil
}
