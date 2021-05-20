package cmd

import (
	"context"
	"fmt"
	"net/url"

	"github.com/mittwald/goharbor-client/v3/apiv1"
	"github.com/spf13/cobra"
	"github.com/x1nchen/portainer-cli/model"
)

var (
	RegistryUser     string
	RegistryPassword string
	RegistryServerAddress string
)

func init() {
	loginRegistryCmd.Flags().StringVarP(&RegistryUser, "user", "u", "", "registry user")
	loginRegistryCmd.Flags().StringVarP(&RegistryPassword, "password", "p", "", "registry password")
	loginRegistryCmd.Flags().StringVarP(&RegistryServerAddress, "serveraddress", "", "", "registry server address")
}

// 登录
var loginRegistryCmd = &cobra.Command{
	Use:     "login-registry",
	Short:   "login registry",
	PreRunE: initUnauthorizedManager,
	Long:    ``,
	RunE:    loginRegistry,
}

func loginRegistry(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	registryClient, err := apiv1.NewRESTClientForHost(
		RegistryServerAddress + "/api", // TODO need carefully handle the trailing slash
		RegistryUser,
		RegistryPassword,
	)

	if err != nil {
		return err
	}

	// try list project to verify the availability of credentials
	// TODO another way to verify the availability of credentials? just like mysql ping?
	_, err = registryClient.ListProjects(ctx, "")
	if err != nil {
		return fmt.Errorf("auth registry error %w", err)
	}

	registryURL, err := url.Parse(RegistryServerAddress)
	if err != nil {
		return err
	}

	// don't need email but used to encode x-registry-auth base64 token
	user := model.RegistryUser{
		Username:      RegistryUser,
		Password:      RegistryPassword,
		Email:         "",
		ServerAddress: registryURL.Host,
	}

	if err = manager.store.RegistryService.UpdateUser(&user); err != nil {
		return err
	}

	cmd.Printf("login registry success\n")
	return err
}
