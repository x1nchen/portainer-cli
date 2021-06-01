package cmd

import (
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
)

func init() {}

// upgrade
var upgradeCmd = &cobra.Command{
	Use:     "upgrade ",
	Short:   "init with terminal box to do deploy interactively",
	PreRunE: initAuthorizedManager,
	Long:    ``,
	RunE:    upgrade,
}

// upgrade single service instance
func upgrade(cmd *cobra.Command, args []string) error {
	// first find out all containers
	containers, err := manager.store.ContainerService.FindAllContainers()
	if err != nil {
		cmd.PrintErr("find containers err", err)
		return err
	}

	var completer = func(d prompt.Document) []prompt.Suggest {
		var s []prompt.Suggest
		for _, container := range containers {
			serviceName := container.Names[0][1:]
			imageTag := strings.SplitN(container.Image, "/", 2)[1]
			description := fmt.Sprintf("%s %s %s", container.EndpointName, imageTag, container.State)
			s = append(s, prompt.Suggest{
				Text:        serviceName,
				Description: description,
			})
		}

		return prompt.FilterContains(s, d.GetWordBeforeCursor(), true)
	}

	cmd.Println("Please select service need to upgrade")
	service := prompt.Input("> ", completer)
	fmt.Println("You selected " + service)
	// select




	return nil
}
