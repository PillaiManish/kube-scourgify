package secrets

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"kube-scourgify/utils"
	"strings"
)

func FindStaleSecrets(ctx context.Context, kubeClient *kubernetes.Clientset, dynamicClient *dynamic.DynamicClient, conditions utils.Conditions, deleteFlagKey bool) error {
	secretsList, err := kubeClient.CoreV1().Secrets("").List(ctx, v1.ListOptions{
		FieldSelector: "type=kubernetes.io/tls",
	})
	if err != nil {
		panic(err)
	}

	var staleSecrets []corev1.Secret

	// finding out the secrets which has cert-manager.io/certificate-name annotations
	for _, secret := range secretsList.Items {
		if secret.OwnerReferences != nil && len(secret.OwnerReferences) > 0 {
			continue
		}
		if secret.Annotations != nil {
			if certName, found := secret.Annotations["cert-manager.io/certificate-name"]; found && utils.IsNotEmpty(certName) {

				certModel, err := utils.GetCertificate(ctx, dynamicClient, secret.Namespace, certName)
				if err != nil {

					// append only if error is notFound -- certificate does not exist
					if errors.IsNotFound(err) {
						fmt.Printf("Certificate %s not found\n", certName)
						staleSecrets = append(staleSecrets, secret)
					}
					continue
				}

				// if certificate.Spec.SecretName is not same as secret name,
				// then certificate is pointing to some other latest secret,
				// so current secret can be considered as stale
				if certModel.Spec.SecretName != secret.Name {
					staleSecrets = append(staleSecrets, secret)
				}

			}
		}
	}

	printStaleSecrets(staleSecrets)

	if deleteFlagKey {
		for _, secret := range staleSecrets {
			err = kubeClient.CoreV1().Secrets(secret.Namespace).Delete(ctx, secret.Name, v1.DeleteOptions{})
			if err != nil {
				return err
			}

		}
	}

	return nil

}

func printStaleSecrets(staleSecrets []corev1.Secret) {
	// Determine the max width for each column
	maxNameLen := len("Secret Name")
	maxNamespaceLen := len("Namespace")

	for _, secret := range staleSecrets {
		if len(secret.Name) > maxNameLen {
			maxNameLen = len(secret.Name)
		}
		if len(secret.Namespace) > maxNamespaceLen {
			maxNamespaceLen = len(secret.Namespace)
		}
	}

	// Add padding for aesthetics
	maxNameLen += 2
	maxNamespaceLen += 2

	// Generate the table format dynamically
	line := fmt.Sprintf("+%s+%s+",
		strings.Repeat("-", maxNameLen+2),
		strings.Repeat("-", maxNamespaceLen+2))

	header := fmt.Sprintf("| %-*s | %-*s |", maxNameLen, "Secret Name", maxNamespaceLen, "Namespace")

	// Print table
	fmt.Println(line)
	fmt.Println(header)
	fmt.Println(line)
	for _, secret := range staleSecrets {
		fmt.Printf("| %-*s | %-*s |\n", maxNameLen, secret.Name, maxNamespaceLen, secret.Namespace)
	}
	fmt.Println(line)
}
