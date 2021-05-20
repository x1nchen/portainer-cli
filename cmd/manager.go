package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mittwald/goharbor-client/v3/apiv1"
	"github.com/spf13/cobra"
	climodel "github.com/x1nchen/portainer-cli/model"

	perr "github.com/pkg/errors"

	"github.com/x1nchen/portainer-cli/cache"
	"github.com/x1nchen/portainer-cli/client"
)

func initManager(store *cache.Store, pclient *client.PortainerClient, cmd *cobra.Command) *Manager {
	m := &Manager{
		store:   store,
		pclient: pclient,
		cmd: cmd,
	}
	return m
}

type Manager struct {
	store          *cache.Store
	cmd            *cobra.Command
	pclient        *client.PortainerClient
	registryClient *apiv1.RESTClient
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

// portainer 服务器数据同步到本地 db 缓存
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

	err = c.store.ContainerService.BatchUpdateContainers(containerList...)
	if err != nil {
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
