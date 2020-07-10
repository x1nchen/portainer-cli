package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	climodel "github.com/x1nchen/portainer-cli/model"

	perr "github.com/pkg/errors"

	"github.com/x1nchen/portainer-cli/cache"
	"github.com/x1nchen/portainer-cli/client"
)

func initManager(store *cache.Store, pclient *client.PortainerClient) *Manager {
	m := &Manager{
		store:   store,
		pclient: pclient,
	}
	return m
}

type Manager struct {
	store   *cache.Store
	pclient *client.PortainerClient
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
			return err
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
		time.Sleep(500 * time.Millisecond)
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

	_, err = c.store.ContainerService.CreateDatabase()
	if err != nil {
		return err
	}

	err = c.store.ContainerService.BatchUpdateContainers(containerList...)
	if err != nil {
		return err
	}

	return nil
}
