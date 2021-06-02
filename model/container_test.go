package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1nchen/portainer/model"
)

func TestContainers(t *testing.T) {
	t.Run("func ContainerName", func(t *testing.T) {
		ce := ContainerExtend{
			EndpointId:   0,
			EndpointName: "",
			DockerContainer: model.DockerContainer{
				ID:    "",
				Names: []string{"/go-report-grpc-srv"},
			},
		}
		assert.Equal(t, "go-report-grpc-srv", ce.ContainerName())
	})
}
