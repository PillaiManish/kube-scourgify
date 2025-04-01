package controller

import (
	"context"
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"kube-scourgify/pkg/controller/certificateRequests"
	"kube-scourgify/pkg/controller/certificates"
	"kube-scourgify/pkg/controller/challenges"
	"kube-scourgify/pkg/controller/orders"
	"kube-scourgify/pkg/controller/secrets"
	"kube-scourgify/utils"
)

func FindStaleResource(kind, group, version, name string) error {
	kubeClient, err := utils.CreateKubeClient()
	if err != nil {
		return err
	}

	dynamicClient, err := utils.CreateDynamicClient()
	if err != nil {
		return err
	}

	ctx := context.TODO()

	discoveryClient := kubeClient.Discovery()

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

	gvr := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: "configmaps",
	}

	resources, err := dynamicClient.Resource(gvr).Namespace("").List(ctx, v1.ListOptions{})
	if err != nil {
		return err
	}

	return findStaleResource(resources, kind)
}

func findStaleResource(resources *unstructured.UnstructuredList, kind string) error {
	switch kind {
	case "secrets":
		return secrets.FindStaleSecrets()
	case "certificates":
		return certificates.FindStaleCertificates()
	case "certificaterequests":
		return certificateRequests.FindStaleCertificateRequests()
	case "orders":
		return orders.FindStaleOrders()
	case "challenges":
		return challenges.FindStaleChallenges()
	default:
		return fmt.Errorf("to be implemented")
	}

}
