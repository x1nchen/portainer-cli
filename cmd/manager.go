package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types/network"
	"github.com/spf13/cobra"
	"github.com/x1nchen/portainer-cli/dockerhub"
	clierr "github.com/x1nchen/portainer-cli/err"
	climodel "github.com/x1nchen/portainer-cli/model"
	"github.com/x1nchen/portainer/model"

	perr "github.com/pkg/errors"

	"github.com/x1nchen/portainer-cli/cache"
	"github.com/x1nchen/portainer-cli/client"
)

func initManager(
	store *cache.Store,
	pclient *client.PortainerClient,
	registryClient *dockerhub.Client,
	cmd *cobra.Command,
) *Manager {
	m := &Manager{
		store:          store,
		pclient:        pclient,
		registryClient: registryClient,
		cmd:            cmd,
	}
	return m
}

type Manager struct {
	store          *cache.Store
	cmd            *cobra.Command
	pclient        *client.PortainerClient
	registryClient *dockerhub.Client
}

func (c *Manager) Login(user string, password string) error {
	if c.pclient == nil {
		return errors.New("pclient not initiated")
	}
	token, err := c.pclient.Auth(context.TODO(), user, password)

	if err != nil {
		return perr.WithMessage(err, "login failed")
	}

	// TODO 登录成功后，将 token 写入缓存
	if err = c.store.TokenService.SaveToken(token); err != nil {
		return perr.WithMessage(err, "save token failed")
	}

	if Verbose {
		c.cmd.Println("[portainer] token", token)
	}
	return nil
}

// SyncData portainer 服务器数据同步到本地 db 缓存
func (c *Manager) SyncData() error {
	ctx := context.Background()

	if c.pclient == nil {
		return errors.New("pclient not initiated")
	}
	eps, err := c.pclient.ListEndpoint(ctx)
	if err != nil {
		return err
	}
	//
	containerList := make([]climodel.ContainerExtend, 0, 200)
	// traverse all endpoints
	// 1. get the container in current endpoint
	// 2. add current endpoint to batch

	updateTime := time.Now()
	for _, ep := range eps {
		cons, err := c.pclient.ListContainer(ctx, int(ep.Id))
		if err != nil {
			fmt.Printf("list container [%s] error %v\n", ep.Name, err)
		}

		for _, con := range cons {
			containerList = append(containerList, climodel.ContainerExtend{
				EndpointId:      int(ep.Id),
				EndpointName:    ep.Name,
				DockerContainer: con,
				UpdateTime:      updateTime,
			})
		}
		// console log
		fmt.Printf("sync endpoint %s container number %d\n", ep.Name, len(cons))

		// force interval to avoid 502 error (api rate limit)
		time.Sleep(200 * time.Millisecond)
	}
	err = c.store.EndpointService.TruncateDatabase()
	if err != nil {
		return err
	}

	_, err = c.store.EndpointService.CreateDatabase()
	if err != nil {
		return err
	}

	// store endpoints
	err = c.store.EndpointService.BatchUpdateEndpoints(eps...)
	if err != nil {
		return err
	}

	err = c.store.ContainerService.TruncateDatabase()
	if err != nil {
		return err
	}

	_, _, err = c.store.ContainerService.CreateDatabase()
	if err != nil {
		return err
	}

	err = c.store.ContainerService.BatchUpdateContainers(context.TODO(), containerList...)
	if err != nil {
		return err
	}

	return nil
}

// UpgradeService Do deploy specified docker with given image tag
func (c *Manager) UpgradeService(
	ctx context.Context,
	containerID string,
	targetImageTag string,
) error {
	oldContainer, err := c.store.ContainerService.GetContainByID(containerID)
	if err != nil {
		return fmt.Errorf("find container failed %w", err)
	}

	containerDetail, _, err := c.pclient.PClient.DockerApi.InspectContainer(
		ctx,
		int32(oldContainer.EndpointId),
		containerID,
	)

	if err != nil {
		return fmt.Errorf("inspect container failed %w", err)
	}

	// 2 create image
	targetImageShortName, targetTag := SplitFullImageName(targetImageTag)
	targetFullImageName := c.registryClient.ServerAddr() + "/" + targetImageTag

	// get registry auth token
	user, err := manager.store.RegistryService.GetUser()
	if err != nil {
		if err == clierr.ErrObjectNotFound {
			return err
		}
	}
	data, _ := json.Marshal(user)
	registryAuthToken := base64.StdEncoding.EncodeToString(data)

	c.cmd.Println("target image name", targetFullImageName)
	_, err = c.pclient.PClient.DockerApi.CreateImage(
		ctx,
		registryAuthToken,
		int32(oldContainer.EndpointId),
		c.registryClient.ServerAddr()+"/"+targetImageShortName,
		targetTag,
	)

	if err != nil {
		c.cmd.PrintErrf("create image %s failed %v\n", targetFullImageName, err)
		return err
	}

	// 3 delete previous container
	_, err = c.pclient.PClient.DockerApi.DeleteContainer(
		ctx,
		int32(oldContainer.EndpointId),
		containerID,
	)

	if err != nil {
		return err
	}

	c.cmd.Println("delete container success ", containerID)

	// change container image
	containerDetail.Config.Image = targetFullImageName
	containerConfig := model.ContainerConfigWrapper{
		Config:     containerDetail.Config,
		HostConfig: containerDetail.HostConfig,
		NetworkingConfig: &network.NetworkingConfig{
			EndpointsConfig: containerDetail.NetworkSettings.Networks},
	}

	// 4 create container
	newContainer, _, err := manager.pclient.PClient.DockerApi.CreateContainer(
		ctx,
		int32(oldContainer.EndpointId),
		containerDetail.Name,
		containerConfig)

	if err != nil {
		return err
	}

	// 5. start container
	c.cmd.Println("create container success", newContainer.ID)
	// TODO we should save the id into our cache store

	_, err = manager.pclient.PClient.DockerApi.StartContainer(ctx,
		int32(oldContainer.EndpointId),
		newContainer.ID)

	if err != nil {
		return err
	}

	// sync container of current endpoint
	// TODO maybe we can do this more frequently
	c.cmd.Println("start container success", newContainer.ID)
	c.cmd.Println("sync endpoint container", oldContainer.EndpointName)

	cons, err := manager.pclient.ListContainer(ctx, oldContainer.EndpointId)
	if err != nil {
		c.cmd.PrintErr(err)
		return err
	}

	updateTime := time.Now()
	endpointContainerList := make([]climodel.ContainerExtend, 0, len(cons))
	for _, con := range cons {
		endpointContainerList = append(endpointContainerList, climodel.ContainerExtend{
			EndpointId:      oldContainer.EndpointId,
			EndpointName:    oldContainer.EndpointName,
			DockerContainer: con,
			UpdateTime:      updateTime,
		})
	}

	if err = manager.store.ContainerService.SyncEndpointContainer(
		ctx,
		oldContainer.EndpointId,
		endpointContainerList...,
	); err != nil {
		return err
	}

	return nil

}

func SplitFullImageName(name string) (imageName, imageTag string) {
	image := strings.Split(name, ":")
	imageName = image[0]
	if len(image) == 1 {
		imageTag = "<none>"
		return
	}
	imageTag = image[1]
	return
}

// SplitFullRegistryImageName
// split dockerRegistryHost/imageShortName/imageTag into each piece
// example:
// dockerhub.com/ccx/go-test-grpc-srv:v1.0.0
// yields [dockerhub.com ccx/go-test-grpc-srv v1.0.0]
func SplitFullRegistryImageName(name string) (dockerRegistryHost, imageShortName, imageTag string) {
	image := strings.Split(name, ":")
	imageName := image[0]
	if len(image) == 1 {
		imageTag = "<none>"
		return
	} else {
		imageTag = image[1]
	}

	imageNameParts := strings.SplitN(imageName, "/", 2)
	if len(imageNameParts) == 1 {
		dockerRegistryHost = imageNameParts[0]
		return
	}
	dockerRegistryHost = imageNameParts[0]
	imageShortName = imageNameParts[1]
	return
}
