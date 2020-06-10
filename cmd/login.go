package cmd

import (
	"context"
	"fmt"
	"github.com/eleztian/portainer"
	"github.com/eleztian/portainer/model"
	"github.com/spf13/cobra"
)

var (
	User string
	Password string
)

func init() {
	loginCmd.Flags().StringVarP(&User, "user", "u", "", "user")
	loginCmd.Flags().StringVarP(&Password, "password", "p", "", "user")
}

// 登录
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login to get the auth token",
	// TODO
	Long:  ``,
	RunE: controller.Login,
}


func login(cmd *cobra.Command, args []string) {
	// TODO module abstraction
	portainerCfg := &portainer.Configuration{
		BasePath:      fmt.Sprintf("%s/api", Host),
	}
	pclient := portainer.NewAPIClient(portainerCfg)
	req := model.AuthenticateUserRequest{
		Username: User,
		Password: Password,
	}
	res, _, err := pclient.AuthApi.AuthenticateUser(context.TODO(), req)
	if err != nil {
		cmd.PrintErr("login failed", err)
		return
	}
	// TODO jwt token store into boltdb
	cmd.Printf("login success, jwt token %s\n", res.Jwt)
}