package cmd

import (
	"fmt"
	"github.com/spf13/cobra"

	"os"
)

var (
	// portainer host
	Host string
)


func init() {
	// cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&Host, "host", "", "host base url such as http://localhost:9000")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(loginCmd)
}

var rootCmd = &cobra.Command{
	Use:   "portainer-cli",
	Short: "portainer-cli is a tool for interacte with portainer",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
