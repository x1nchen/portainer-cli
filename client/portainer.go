package client

import (
	"context"
	"errors"
	"fmt"
	"strings"

	perr "github.com/pkg/errors"

	"github.com/x1nchen/portainer"
	"github.com/x1nchen/portainer/model"
)

var (
	ErrorAuthFailed = errors.New("auth failed")
)

type PortainerClient struct {
	PClient *portainer.APIClient
}

// TODO options style DI
func NewPortainerClient(host string, jwtToken string) *PortainerClient {
	portainerCfg := &portainer.Configuration{
		BasePath: fmt.Sprintf("%s/api", host),
	}

	c := portainer.NewAPIClient(portainerCfg)
	c.JwtToken = jwtToken

	return &PortainerClient{
		PClient: c,
	}
}

// Auth no need with jwt token
func (p *PortainerClient) Auth(ctx context.Context, user string, password string) (string, error) {
	// fmt.Println(user, password)
	req := model.AuthenticateUserRequest{
		Username: user,
		Password: password,
	}
	authRes, _, err := p.PClient.AuthApi.AuthenticateUser(ctx, req)

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "invalid credentials") {
			return "", ErrorAuthFailed
		}
		return "", perr.WithMessage(err, "auth failed")
	}
	return authRes.Jwt, nil
}

// ListContainer jwt token needed
func (p *PortainerClient) ListContainer(ctx context.Context, endpointId int) (model.DockerContainerListResponse, error) {
	res, _, err := p.PClient.DockerApi.ListContainer(ctx, int32(endpointId))

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "invalid credentials") {
			return res, ErrorAuthFailed
		}
	}
	return res, nil
}
