package controller

import (
	"fmt"
	"kube-scourgify/utils"
)

func FindStaleResource(kind, group, version, name string) error {
	kubeClient, err := utils.CreateKubeClient()
	if err != nil {
		return err
	}

	discoveryClient := kubeClient.Discovery()

	// Check if it's a built-in resource
	resourceList, err := discoveryClient.ServerResourcesForGroupVersion(version)
	if err != nil {
		return err
	}

	fmt.Println("chcel", resourceList.APIResources)
	return nil
}
