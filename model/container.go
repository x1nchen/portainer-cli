package model

import "github.com/x1nchen/portainer/model"

type ContainerExtend struct {
	EndpointId   int    `json:"endpoint_id"`
	EndpointName string `json:"endpoit_name"`
	model.DockerContainer
}
