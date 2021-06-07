package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
	ctx := context.Background()
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

	var instanceAnswers []string

	// select instance
	{
		var serviceInstanceOptions []string
		for _, container := range containers {
			srvName := container.ContainerName()
			if srvName == serviceName {
				option := fmt.Sprintf("%s %s %s %s", container.ID, container.EndpointName, container.Image, container.State)
				serviceInstanceOptions = append(serviceInstanceOptions, option)
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

		err = survey.Ask(serviceInstanceQuestion, &instanceAnswers)
		if err != nil {
			return err
		}
		cmd.Println("You selected ", instanceAnswers)
	}

	if len(instanceAnswers) == 0 {
		// cmd.Println("You selected ", instanceAnswers)
		return errors.New("at lease one instance must be selected")
	}

	instanceAnswer := instanceAnswers[0]
	_, _, imageName, _ := GetContainerSpec(instanceAnswer)

	_, imageShortName, _ := SplitFullRegistryImageName(imageName)

	tagList, err := manager.registryClient.FindImageTagList(ctx, imageShortName)
	if err != nil {
		return err
	}

	var imageTagAnswer string

	// select image tag
	{
		var imageTagOptions []string
		for _, tagModel := range tagList {
			option := fmt.Sprintf("%s:%s", imageShortName, tagModel.Name)
			imageTagOptions = append(imageTagOptions, option)
		}

		// the questions to ask
		var imageTagQuestion = []*survey.Question{
			{
				Name: "imagetag",
				Prompt: &survey.Select{
					Message: "select instance to upgrade>",
					Options: imageTagOptions,
				},
				Validate: survey.Required,
			},
		}

		err = survey.Ask(imageTagQuestion, &imageTagAnswer)
		if err != nil {
			return err
		}
		cmd.Println("You selected ", imageTagAnswer)
	}

	for _, instanceAnswer := range instanceAnswers {
		instanceID := strings.Split(instanceAnswer, " ")[0]
		container, err := manager.store.ContainerService.GetContainByID(instanceID)
		if err != nil {
			cmd.Printf("find service container error %s \n", instanceAnswer)
			continue
		}
		if err = manager.UpgradeService(ctx, container.ID, imageTagAnswer); err != nil {
			cmd.Printf("upgrade service container error %s \n", instanceAnswer)
			continue
		}
		cmd.Printf("[>] %s %s %s\n", container.ContainerName(), container.EndpointName, container.Image)
	}

	return nil
}

func GetContainerSpec(instanceAnswer string) (containerID, endpointName, imageName, state string) {
	specs := strings.Split(instanceAnswer, " ")
	switch len(specs) {
	case 0:
		return
	case 1:
		containerID = specs[0]
		return
	case 2:
		containerID, endpointName = specs[0], specs[1]
		return
	case 3:
		containerID, endpointName, imageName = specs[0], specs[1], specs[2]
		return
	case 4:
		containerID, endpointName, imageName, state = specs[0], specs[1], specs[2], specs[3]
		return
	}
	return
}
