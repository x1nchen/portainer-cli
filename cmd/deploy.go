package cmd

import (
	"context"
	"errors"

	"github.com/spf13/cobra"
)

func init() {}

// deploy
var deployCmd = &cobra.Command{
	Use:     "deploy",
	Short:   "deploy container with given image tag",
	PreRunE: initAuthorizedManager,
	Long:    ``,
	RunE:    deploy,
}

// STEP
// 1. inspect a container
// 2. create deploy target image (docker pull)
// 3. delete target container
// 4. create container
// 5. start container
func deploy(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	name := args[0]
	targetImageTag := args[1]

	containers, err := manager.store.ContainerService.FuzzyFindContainerByName(name)
	if err != nil {
		cmd.PrintErr(err)
		return err
	}
	if len(containers) == 0 {
		cmd.PrintErrf("%s service not found\n", name)
		return errors.New("service not found")
	}

	// TODO should sync the container in current endpoint instance before try upgrade
	if len(containers) > 1 {
		cmd.PrintErrf("%s service has more than 1 instance\n", name)
		return errors.New("service has more than one instance")
	}

	container := containers[0]
	err = manager.UpgradeService(ctx, container.ID, targetImageTag)
	if err != nil {
		return err
	}

	return nil
}
