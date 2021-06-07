module github.com/x1nchen/portainer-cli

// replace github.com/x1nchen/portainer => ../portainer
replace github.com/mittwald/goharbor-client/v3 => github.com/x1nchen/goharbor-client/v3 v3.3.1-0.20210607124842-d4cf1cb80f59

go 1.16

require (
	github.com/AlecAivazis/survey/v2 v2.2.12
	github.com/docker/docker v20.10.6+incompatible
	github.com/json-iterator/go v1.1.10
	github.com/logrusorgru/aurora v2.0.3+incompatible
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/mittwald/goharbor-client/v3 v3.3.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.7.0
	github.com/x1nchen/portainer v1.23.8
	go.etcd.io/bbolt v1.3.5
)

// lint error
retract v0.0.7
