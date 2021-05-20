package model

import (
	"fmt"

	"github.com/x1nchen/portainer/model"
)

type ContainerExtend struct {
	EndpointId   int    `json:"endpoint_id"`
	EndpointName string `json:"endpoit_name"`
	model.DockerContainer
}

func (c ContainerExtend) KeyWithEndpoint() string {
	var containerName string

	if len(c.Names) > 0 {
		if len(c.Names[0]) > 0 {
			// 注意：容器的名字有前缀 "/"，如 /node-api
			containerName = c.Names[0][1:]
		}
	}

	return fmt.Sprintf("%s#%d", containerName, c.EndpointId)
}

func (c ContainerExtend) KeyWithContainerID() string {
	return c.ID
}

// RegistryUser is dockerhub credentials
type RegistryUser struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	ServerAddress string `json:"serveraddress"`
}