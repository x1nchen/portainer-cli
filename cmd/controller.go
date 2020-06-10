package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/x1nchen/portainer-cli/cache"
	"github.com/x1nchen/portainer-cli/client"
)

var controller *Controller

func initController(store cache.Store, client client.PortainerClient) {
	controller = &Controller{
		store:  store,
		client: client,
	}
}

type Controller struct {
	store cache.Store
	client client.PortainerClient
}

func (c *Controller) Login(cmd *cobra.Command, args []string) error {
	_, err := c.client.Auth(context.TODO(), User, Password)

	if err != nil {
		cmd.PrintErrf("login failed %v", err)
		return err
	}

	// 登录成功后，将 token 写入缓存
	// c.store.
	return nil
}