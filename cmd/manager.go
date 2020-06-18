package cmd

import (
	"context"
	"errors"

	perr "github.com/pkg/errors"

	"github.com/x1nchen/portainer-cli/cache"
	"github.com/x1nchen/portainer-cli/client"
)

func initManager(store cache.Store, pclient *client.PortainerClient) *Manager {
	m := &Manager{
		store:   store,
		pclient: pclient,
	}
	return m
}

type Manager struct {
	store   cache.Store
	pclient *client.PortainerClient
}

func (c *Manager) Login(user string, password string) error {
	if c.pclient == nil {
		return errors.New("pclient not initiated")
	}
	token, err := c.pclient.Auth(context.TODO(), user, password)
	// fmt.Println(token)
	if err != nil {
		return perr.WithMessage(err, "login failed")
	}

	// TODO 登录成功后，将 token 写入缓存
	if err = c.store.SaveToken(token); err != nil {
		return perr.WithMessage(err, "save token failed")
	}

	return nil
}
