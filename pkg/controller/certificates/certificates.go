package certificates

import (
	"context"
	"fmt"
	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"kube-scourgify/utils"
	"strings"
)

func FindStaleCertificates(ctx context.Context, dynamicClient *dynamic.DynamicClient, conditions utils.Conditions) error {
	var staleCertificates []*certmanagerv1.Certificate

	// gvr for certificates
	gvr := certmanagerv1.SchemeGroupVersion.WithResource(utils.CERTIFICATES)

	//	Get list of certificates
	unStructuredCertificates, err := dynamicClient.Resource(gvr).List(context.Background(), v1.ListOptions{
		// TODO: check if any filter needed
	})

	if err != nil {
		return err
	}

	// loop through each certificate
	for _, unStructuredCertificate := range unStructuredCertificates.Items {
		cert := &certmanagerv1.Certificate{}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(unStructuredCertificate.UnstructuredContent(), cert)
		if err != nil {
			return err
		}

		// if the issuerRef.Group == cert-manager.io
		if cert.Spec.IssuerRef.Group == certmanagerv1.SchemeGroupVersion.Group {
			_, _, err = utils.GetIssuerByName(ctx, dynamicClient, cert.Spec.IssuerRef.Kind, cert.Namespace, cert.Spec.IssuerRef.Name)
			if err != nil {
				if errors.IsNotFound(err) {
					staleCertificates = append(staleCertificates, cert)
				}
			}
		}
	}

	printStaleCertificates(staleCertificates)

	return nil
}

func printStaleCertificates(staleCertificates []*certmanagerv1.Certificate) {
	// Determine the max width for each column
	maxNameLen := len("Certificate Name")
	maxNamespaceLen := len("Namespace")

	for _, certificate := range staleCertificates {
		if len(certificate.Name) > maxNameLen {
			maxNameLen = len(certificate.Name)
		}
		if len(certificate.Namespace) > maxNamespaceLen {
			maxNamespaceLen = len(certificate.Namespace)
		}
	}

	// Add padding for aesthetics
	maxNameLen += 2
	maxNamespaceLen += 2

	// Generate the table format dynamically
	line := fmt.Sprintf("+%s+%s+",
		strings.Repeat("-", maxNameLen+2),
		strings.Repeat("-", maxNamespaceLen+2))

	header := fmt.Sprintf("| %-*s | %-*s |", maxNameLen, "Certificate Name", maxNamespaceLen, "Namespace")

	// Print table
	fmt.Println(line)
	fmt.Println(header)
	fmt.Println(line)
	for _, certificate := range staleCertificates {
		fmt.Printf("| %-*s | %-*s |\n", maxNameLen, certificate.Name, maxNamespaceLen, certificate.Namespace)
	}
	fmt.Println(line)
}
