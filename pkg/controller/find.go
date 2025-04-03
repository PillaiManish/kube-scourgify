package controller

import (
	"context"
	"fmt"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"kube-scourgify/pkg/controller/certificateRequests"
	"kube-scourgify/pkg/controller/certificates"
	"kube-scourgify/pkg/controller/challenges"
	"kube-scourgify/pkg/controller/orders"
	"kube-scourgify/pkg/controller/secrets"
	"kube-scourgify/utils"
)

func FindStaleResource(kind, group, version, filepath string) error {
	ctx := context.Background()

	kubeClient, err := utils.CreateKubeClient()
	if err != nil {
		return err
	}

	dynamicClient, err := utils.CreateDynamicClient()
	if err != nil {
		return err
	}

	discoveryClient := kubeClient.Discovery()

	if kind == "" {
		return fmt.Errorf("kind is required")
	}

	if group == "" || version == "" {
		resourceList, err := discoveryClient.ServerResourcesForGroupVersion(version)
		if err != nil {
			return err
		}

		for _, resource := range resourceList.APIResources {
			if resource.Kind == kind {
				group = resource.Group
				version = "v1" //FIXME: check how to handle version & group properly
				break
			}
		}
	}

	conditions, err := utils.ParseConditions(filepath)

	return findStaleResource(ctx, kubeClient, dynamicClient, kind, conditions)
}

func findStaleResource(ctx context.Context, kubeClient *kubernetes.Clientset, dynamicClient *dynamic.DynamicClient, kind string, conditions utils.Conditions) error {
	switch kind {
	case "secrets":
		return secrets.FindStaleSecrets(ctx, kubeClient, dynamicClient, conditions)
	case utils.CERTIFICATES:
		return certificates.FindStaleCertificates(ctx, dynamicClient, conditions)
	case utils.CERTIFICATEREQUESTS:
		return certificateRequests.FindStaleCertificateRequests(ctx, dynamicClient, conditions)
	case utils.ORDERS:
		return orders.FindStaleOrders(ctx, dynamicClient, conditions)
	case utils.CHALLENGES:
		return challenges.FindStaleChallenges(ctx, dynamicClient, conditions)
	default:
		return fmt.Errorf("to be implemented")
	}

}
