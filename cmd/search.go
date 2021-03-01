package cmd

import (
	"fmt"
	"strings"

	climodel "github.com/x1nchen/portainer-cli/model"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

func init() {
	// searchCmd.Flags().StringVarP(&User, "user", "u", "", "user")
	// searchCmd.Flags().StringVarP(&Password, "password", "p", "", "user")
}

// 登录
var searchCmd = &cobra.Command{
	Use:     "search",
	Short:   "search container with fuzzy name",
	PreRunE: initAuthorizedManager,
	Long:    ``,
	RunE:    search,
}

// STEP
// 1. fuzzy match the container name or image name?
// 2. search from cache(boltdb) with step 1
// 3. verify the container from cache by call docker api
// 4. show container name and node name with list formation
func search(cmd *cobra.Command, args []string) error {
	// TODO make color output flagged
	au := aurora.NewAurora(true)

	name := args[0]
	containers, err := manager.store.ContainerService.FuzzyFindContainerByName(name)
	if err != nil {
		cmd.PrintErr(err)
		return err
	}

	title := fmt.Sprintf(genOutTemplate(containers), "NAME", "TAG", "NODE", "STATE", "IMAGE")
	cmd.Println(au.White(title))
	for _, container := range containers {
		image := strings.Split(container.Image, ":")
		outMessage := fmt.Sprintf(genOutTemplate(containers), container.Names[0][1:],
			image[1], container.EndpointName, container.State, image[0])
		if container.State == "running" { // TODO const
			cmd.Println(au.Green(outMessage))
		}
		if container.State == "exited" { // TODO const
			cmd.Println(au.Red(outMessage))
		}

	}

	return nil
}

// get the longest field to generate template
func genOutTemplate(list []climodel.ContainerExtend) string {
	var nameLen, tagLen, nodeLen, stateLen, imageLen int
	for _, v := range list {
		image := strings.Split(v.Image, ":")
		if len(v.Names[0][1:]) > nameLen {
			nameLen = len(v.Names[0][1:])
		}
		if len(image[0]) > imageLen {
			imageLen = len(image[0])
		}
		if len(image[1]) > tagLen {
			tagLen = len(image[1])
		}
		if len(v.EndpointName) > nodeLen {
			nodeLen = len(v.EndpointName)
		}
		if len(v.State) > stateLen {
			stateLen = len(v.State)
		}
	}
	return fmt.Sprint(
		"%-", nameLen, "s ",
		"%-", tagLen, "s ",
		"%-", nodeLen, "s ",
		"%-", stateLen, "s ",
		"%-", imageLen, "s ",
	)
}
