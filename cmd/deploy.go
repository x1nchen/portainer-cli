package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	clierr "github.com/x1nchen/portainer-cli/err"
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
// 1. inspect container
// 2. create target image
// 3. delete target container
// 4. create container
func deploy(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	name := args[0]
	targetImageTag := args[1]

	containers, err := manager.store.ContainerService.FuzzyFindContainerByName(name)
	if err != nil {
		cmd.PrintErr(err)
		return err
	}

	// TODO should sync the container in current endpoint instance before try upgrade
	if len(containers) > 1 {
		cmd.PrintErrf("%s service has more than 1 instance\n", name)
		return errors.New("service has more than one instance")
	}

	container := containers[0]
	// 1 inspect container
	containerJSON, _, err := manager.pclient.PClient.DockerApi.InspectContainer(
		ctx,
		int32(container.EndpointId),
		container.ID,
	)

	if err != nil {
		return fmt.Errorf("inspect container failed %w", err)
	}

	// 2 create image
	targetImageShortName, targetTag := SplitFullImageName(targetImageTag)
	registryHost, _, _ := SplitFullRegistryImageName(container.Image)
	targetFullImageName := registryHost + "/" + targetImageShortName + ":" + targetTag

	// get registry auth token
	user, err := manager.store.RegistryService.GetUser()
	if err != nil {
		if err == clierr.ErrObjectNotFound {
			cmd.Println("registry credentials not found, please try login-registry command first")
			return err
		}
	}
	fmt.Println(user)

	data, _ := json.Marshal(user)
	registryAuthToken := base64.StdEncoding.EncodeToString(data)
	fmt.Println(registryAuthToken)

	cmd.Println("target image name", targetFullImageName)
	_, err = manager.pclient.PClient.DockerApi.CreateImage(
		ctx,
		registryAuthToken,
		int32(container.EndpointId),
		registryHost+"/"+targetImageShortName,
		targetTag,
	)

	if err != nil {
		cmd.PrintErrf("create image %s failed %v\n", targetFullImageName, err)
		return err
	}

	// 3 delete previous container
	_, err = manager.pclient.PClient.DockerApi.DeleteContainer(
		ctx,
		int32(container.EndpointId),
		container.ID,
	)

	if err != nil {
		cmd.PrintErr(err)
		return err
	}
	cmd.Println("delete container success ", container.ID)

	containerJSON.Image = targetFullImageName
	// 4 create container
	newContainer, _, err := manager.pclient.PClient.DockerApi.CreateContainer(
		ctx,
		int32(container.EndpointId),
		containerJSON.Name,
		containerJSON)

	if err != nil {
		cmd.PrintErr(err)
		return err
	}

	// 5. start container
	cmd.Println("create container success", newContainer.ID)
	// TODO we should save the id into our cache store

	_, err = manager.pclient.PClient.DockerApi.StartContainer(ctx,
		int32(container.EndpointId),
		newContainer.ID)

	if err != nil {
		cmd.PrintErr(err)
		return err
	}
	cmd.Println("start container success", newContainer.ID)
	return nil
}
