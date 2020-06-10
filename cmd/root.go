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
	rootCmd.AddCommand(configCmd)
}

var rootCmd = &cobra.Command{
	Use:   "portainer-cli",
	Short: "Portainer CLI",
	Long: `Work seamlessly with Portainer from the command line.`,
	PreRunE: prepare,
	SilenceErrors: true,
	SilenceUsage:  true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func prepare(cmd *cobra.Command, args []string) error {
	return nil
}

