module github.com/x1nchen/portainer-cli

// replace github.com/x1nchen/portainer => ../portainer

go 1.16

require (
	github.com/docker/docker v20.10.6+incompatible // indirect
	github.com/json-iterator/go v1.1.10
	github.com/logrusorgru/aurora v2.0.3+incompatible
	github.com/mittwald/goharbor-client/v3 v3.3.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.7.0
	github.com/x1nchen/portainer v1.23.8
	go.etcd.io/bbolt v1.3.5
)

// lint error
retract v0.0.7
