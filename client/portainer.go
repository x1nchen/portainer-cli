package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/eleztian/portainer"
	"github.com/eleztian/portainer/model"
	"strings"
)

var (
	ErrorAuthFailed = errors.New("auth failed")

)

type PortainerClient struct {
	PClient *portainer.APIClient
}

func NewPortainerClient(host string) *PortainerClient {
	portainerCfg := &portainer.Configuration{
		BasePath:      fmt.Sprintf("%s/api", host),
	}

	c := portainer.NewAPIClient(portainerCfg)

	return &PortainerClient{
		PClient: c,
	}
}

func (p *PortainerClient) Auth(ctx context.Context, user string, password string) (string, error) {
	req := model.AuthenticateUserRequest{
		Username: user,
		Password: password,
	}
	authRes, _, err := p.PClient.AuthApi.AuthenticateUser(ctx, req)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "invalid credentials") {
			return "", ErrorAuthFailed
		}
	}
	return authRes.Jwt, nil
}
