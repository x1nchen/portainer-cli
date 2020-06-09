package cmd

import (
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
	Run: login,
}

func login(cmd *cobra.Command, args []string) {
	
}