package dockerhub

import (
	"context"
	"fmt"
	"strings"

	"github.com/mittwald/goharbor-client/v3/apiv1"
	goharbormodel "github.com/mittwald/goharbor-client/v3/apiv1/model"
)

type Client struct {
	// serverAddr is for concatenate the api path for request
	serverAddr string
	user       string
	password   string
	restClient *apiv1.RESTClient
}

// ServerAddr server address
func (c *Client) ServerAddr() string {
	return c.serverAddr
}

func NewClient(serverAddr, user, password string) (*Client, error) {
	var endpoint = serverAddr
	if !strings.HasPrefix(serverAddr, "https") &&
		!strings.HasPrefix(serverAddr, "http") {
		endpoint = "https://" + serverAddr + "/api"// TODO need carefully handle the trailing slash
	}

	registryClient, err := apiv1.NewRESTClientForHost(
		// serverAddr+"/api",
		endpoint,
		user,
		password,
	)

	if err != nil {
		return nil, err
	}

	c := &Client{
		serverAddr: serverAddr,
		user:       user,
		password:   password,
		restClient: registryClient,
	}

	return c, nil
}

// Auth dockerhub
// TODO another way to verify the availability of credentials? just like mysql ping?
func (c *Client) Auth(ctx context.Context) error {
	_, err := c.restClient.ListProjects(ctx, "")
	if err != nil {
		return fmt.Errorf("docker registry auth error %v", err)
	}

	return nil
}

func (c *Client) FindImageTagList(ctx context.Context, imageShortName string) ([]*goharbormodel.DetailedTag, error) {
	tags, err := c.restClient.GetRepositoryTags(ctx, imageShortName)
	if err != nil {
		return nil, err
	}

	return tags, nil
}
