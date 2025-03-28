package utils

import (
	"os"
	"strings"
)

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
