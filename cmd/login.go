package cmd

import (
	"github.com/spf13/cobra"
)

var (
	User     string
	Password string
)

func init() {
	loginCmd.Flags().StringVarP(&User, "user", "u", "", "user")
	loginCmd.Flags().StringVarP(&Password, "password", "p", "", "user")
}

// 登录
var loginCmd = &cobra.Command{
	Use:     "login",
	Short:   "login to get the auth token",
	PreRunE: initUnauthorizedManager,
	Long:    ``,
	RunE:    login,
}

func login(cmd *cobra.Command, args []string) error {
	if Verbose {
		cmd.Println("[bbolt] db name", manager.store.DBName)
	}
	err := manager.Login(User, Password)
	if err != nil {
		cmd.PrintErr("login failed", err)
		return err
	}

	cmd.Printf("login success\n")
	return err
}
