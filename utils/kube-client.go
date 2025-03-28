package utils

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"fmt"
)

func CreateKubeClient() (*kubernetes.Clientset, error) {
	var kubeConfigPath = getKubeConfigPath()
	var kubeConfig *rest.Config

	if kubeConfigPath != "" {
		config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			return nil, fmt.Errorf("unable to load kubeconfig from %s: %v", kubeConfigPath, err)
		}
		kubeConfig = config
	} else {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("unable to load in-cluster config: %v", err)
		}
		kubeConfig = config
	}

	client, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create a client: %v", err)
	}

	return client, nil
}

func CreateDynamicClient() (*dynamic.DynamicClient, error) {
	var kubeConfigPath = getKubeConfigPath()
	var kubeConfig *rest.Config
	if kubeConfigPath != "" {
		config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			return nil, fmt.Errorf("unable to load kubeconfig from %s: %v", kubeConfigPath, err)
		}
		kubeConfig = config
	} else {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("unable to load in-cluster config: %v", err)
		}
		kubeConfig = config
	}
	client, err := dynamic.NewForConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create dynamic client: %v", err)
	}
	return client, nil
}
