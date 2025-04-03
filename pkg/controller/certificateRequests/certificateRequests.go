package certificateRequests

import (
	"context"
	"fmt"
	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"kube-scourgify/utils"
	"strings"
)

func FindStaleCertificateRequests(ctx context.Context, dynamicClient *dynamic.DynamicClient, conditions utils.Conditions, deleteFlag bool) error {
	var isCertificateDepFlagEnabled, isIssuerDepFlagEnabled bool

	for _, i := range conditions.Deps {
		if strings.ToLower(i) == utils.CERTIFICATE || strings.ToLower(i) == utils.CERTIFICATES {
			isCertificateDepFlagEnabled = true
		} else if strings.ToLower(i) == utils.ISSUER || strings.ToLower(i) == utils.ISSUERS {
			isIssuerDepFlagEnabled = true
		}
	}

	var staleCertificateRequests []*certmanagerv1.CertificateRequest
	gvr := certmanagerv1.SchemeGroupVersion.WithResource(utils.CERTIFICATEREQUESTS)

	unStructuredCertificateRequests, err := dynamicClient.Resource(gvr).List(ctx, v1.ListOptions{})
	if err != nil {
		return err
	}

	for _, unStructuredCertificateRequest := range unStructuredCertificateRequests.Items {
		certRequest := &certmanagerv1.CertificateRequest{}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(unStructuredCertificateRequest.UnstructuredContent(), certRequest)
		if err != nil {
			return err
		}

		isIssuerRefPresent, isCertificatePresent := false, false

		// isIssuerDepFlagEnabled is enabled check for missing issuer's CertificateRequests
		if isIssuerDepFlagEnabled {
			if certRequest.Spec.IssuerRef.Group == certmanagerv1.SchemeGroupVersion.Group {
				_, _, err = utils.GetIssuerByName(ctx, dynamicClient, certRequest.Spec.IssuerRef.Kind, certRequest.Namespace, certRequest.Spec.IssuerRef.Name)

				if err == nil {
					isIssuerRefPresent = true
				} else {
					fmt.Printf("certificate request: %s does not have issuers ref: %s \n", certRequest.Name, certRequest.Spec.IssuerRef.Name)
				}
			}
		}

		if isCertificateDepFlagEnabled {
			if certRequest.Annotations != nil {

				if certificateName, found := certRequest.Annotations[certmanagerv1.CertificateNameKey]; found && certificateName != "" {
					_, err = utils.GetCertificate(ctx, dynamicClient, certRequest.Namespace, certificateName)
					if err == nil {
						isCertificatePresent = true
					}
				} else {
					fmt.Printf("certificate request: %s does not have certificate: %s \n", certRequest.Name, certificateName)

				}
			}
		}

		if isCertificatePresent || isIssuerRefPresent {
			staleCertificateRequests = append(staleCertificateRequests, certRequest)
		}

	}

	fmt.Println(staleCertificateRequests)

	if deleteFlag {
		for _, certRequest := range staleCertificateRequests {
			err = dynamicClient.Resource(gvr).Delete(ctx, certRequest.Name, v1.DeleteOptions{})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
