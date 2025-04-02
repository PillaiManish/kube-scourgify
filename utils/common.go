package utils

import (
	"context"
	"encoding/json"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"os"
	"strconv"
	"strings"
	"time"
)

type Conditions struct {
	Deps []string `json:"deps"`
}

func getKubeConfigPath() string {
	kubeConfig := os.Getenv("KUBECONFIG")
	if isEmpty(kubeConfig) {
		kubeConfig = DEFAULT_CONFIG_PATH
	}
	return kubeConfig
}

func isEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func isNotEmpty(s string) bool {
	return !isEmpty(s)
}

func ParseConditions(filePath string) (Conditions, error) {
	byteValues, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
	}

	var conditions Conditions
	json.Unmarshal(byteValues, &conditions)
	fmt.Println(conditions.Deps)
	return conditions, nil
}

func CreateStaleResourceCR(resourceGroup, resourceType string, conditions, resourceList []string, client dynamic.DynamicClient) error {
	gvr := schema.GroupVersionResource{
		Group:    "kubescourgify.io",
		Version:  "v1",
		Resource: "staleresources",
	}

	obj := map[string]interface{}{
		"apiVersion": "",
		"kind":       "StaleResource",
		"metadata": map[string]interface{}{
			"name": resourceType + strconv.FormatInt(time.Now().Unix(), 10),
		},
		"spec": map[string]interface{}{
			"resourceConditions": conditions,
			"resourceGroup":      resourceGroup,
			"resourceType":       resourceType,
			"resourceList":       resourceList,
		},
	}
	//create a StaleResource CR through dynamic client
	_, err := client.Resource(gvr).Create(context.Background(), &unstructured.Unstructured{Object: obj}, metav1.CreateOptions{})
	return err
}
