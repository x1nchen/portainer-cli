package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

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

// NewPortainerClient
// TODO options style DI
func NewPortainerClient(host string, jwtToken string) *PortainerClient {
	portainerCfg := &portainer.Configuration{
		BasePath:   fmt.Sprintf("%s/api", host),
		HTTPClient: http.DefaultClient,
	}

	// create image should wait for more than 1 minute
	portainerCfg.HTTPClient.Timeout = 60 * time.Second

	c := portainer.NewAPIClient(portainerCfg)
	c.JwtToken = jwtToken

	return &PortainerClient{
		PClient: c,
	}
}

func (p *PortainerClient) CarryToken(token string) {
	p.PClient.JwtToken = token
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
		return res, perr.WithMessage(err, "portainer list container")
	}
	return res, nil
}

// ListEndpoint jwt token needed
func (p *PortainerClient) ListEndpoint(ctx context.Context) (model.EndpointListResponse, error) {
	res, _, err := p.PClient.EndpointsApi.EndpointList(ctx)

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "invalid credentials") {
			return res, ErrorAuthFailed
		}
		return res, perr.WithMessage(err, "portainer list endpoint")
	}

	return res, nil
}

func (p *PortainerClient) GetStatus(ctx context.Context) (model.Status, error) {
	res, _, err := p.PClient.StatusApi.StatusInspect(ctx)

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "invalid credentials") {
			return res, ErrorAuthFailed
		}
		return res, perr.WithMessage(err, "get portainer status")
	}

	return res, nil
}
