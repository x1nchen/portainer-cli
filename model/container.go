package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/x1nchen/portainer/model"
)

const (
	FormatDayLayoutDetail = "2006-01-02 15:04:05"
)

type ContainerExtend struct {
	EndpointId   int    `json:"endpoint_id"`
	EndpointName string `json:"endpoit_name"`
	model.DockerContainer
	// update time for db record
	UpdateTime time.Time `json:"update_time"`
}

// KeyWithEndpoint
// format: {EndpointId}:{ContainerName}:{ContainerID}
func (c ContainerExtend) KeyWithEndpoint() string {
	var containerName string

	if len(c.Names) > 0 {
		if len(c.Names[0]) > 0 {
			// 注意：容器的名字有前缀 "/"，如 /node-api
			containerName = c.Names[0][1:]
		}
	}

	return fmt.Sprintf("%d:%s:%s", c.EndpointId, containerName, c.ID)
}

func (c ContainerExtend) ContainerName() string {
	var containerName string

	if len(c.Names) > 0 {
		if len(c.Names[0]) > 0 {
			// 注意：容器的名字有前缀 "/"，如 /node-api
			name := c.Names[0]
			if strings.HasPrefix(name, "/") {
				containerName = c.Names[0][1:]
			} else {
				containerName = name
			}
		}
	}

	return containerName
}

func (c ContainerExtend) KeyWithContainerID() string {
	return c.ID
}

func (c ContainerExtend) UpdateTimeStr() string {
	return c.UpdateTime.Format(FormatDayLayoutDetail)
}

// RegistryUser is dockerhub credentials
// ref https://docs.docker.com/engine/api/v1.30/#section/Authentication
type RegistryUser struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	ServerAddress string `json:"serveraddress"`
}
