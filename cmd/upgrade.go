package cmd

import (
	"github.com/AlecAivazis/survey/v2"
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

	nameSet := make(map[string]int)
	var serviceNameOptions []string
	for _, container := range containers {
		srvName := container.ContainerName()
		if _, ok := nameSet[srvName]; ok {
			continue
		}

		nameSet[srvName] = 1

		// tmpStrList := strings.SplitN(container.Image, "/", 2)
		// var imageTag string
		// imageTag = tmpStrList[0]
		// if len(tmpStrList) == 2 {
		// 	imageTag = tmpStrList[1]
		// }

		// description := fmt.Sprintf("%s %s %s %s",srvName, container.EndpointName, imageTag, container.State)
		serviceNameOptions = append(serviceNameOptions, srvName)
	}

	// select service name
	var serviceName string
	// the questions to ask
	var serviceNameQuestion = []*survey.Question{
		{
			Name: "servicename",
			Prompt: &survey.Select{
				Message: "select service need to upgrade >",
				Options: serviceNameOptions,
			},
			Validate: survey.Required,
		},
	}

	err = survey.Ask(serviceNameQuestion, &serviceName)
	if err != nil {
		return err
	}

	cmd.Println("You selected " + serviceName)

	// select instance
	{
		var instance string
		var serviceInstanceOptions []string
		for _, container := range containers {
			srvName := container.ContainerName()
			if srvName == serviceName {
				serviceInstanceOptions = append(serviceInstanceOptions, container.ID)
			}
		}

		// the questions to ask
		var serviceInstanceQuestion = []*survey.Question{
			{
				Name: "serviceinstance",
				Prompt: &survey.MultiSelect{
					Message: "select instance to upgrade>",
					Options: serviceInstanceOptions,
				},
				Validate: survey.Required,
			},
		}

		err = survey.Ask(serviceInstanceQuestion, &instance)
		if err != nil {
			return err
		}
		cmd.Println("You selected " + instance)
	}
	return nil
}
