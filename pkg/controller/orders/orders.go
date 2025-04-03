package orders

import (
	"context"
	acmev1 "github.com/cert-manager/cert-manager/pkg/apis/acme/v1"
	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"kube-scourgify/utils"
)

func FindStaleOrders(ctx context.Context, dynamicClient *dynamic.DynamicClient, conditions utils.Conditions) error {
	var staleOrders []*acmev1.Order
	gvr := certmanagerv1.SchemeGroupVersion.WithResource(utils.ORDERS)

	unstructuredOrders, err := dynamicClient.Resource(gvr).List(ctx, v1.ListOptions{})
	if err != nil {
		return err
	}

	for _, unStructuredOrder := range unstructuredOrders.Items {
		order := &acmev1.Order{}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(unStructuredOrder.UnstructuredContent(), order)
		if err != nil {
			return err
		}

		var isIssuerRefPresent, isCertificatePresent bool
		if order.Spec.IssuerRef.Group == acmev1.SchemeGroupVersion.Group {
			_, _, err = utils.GetIssuerByName(ctx, dynamicClient, order.Spec.IssuerRef.Kind, order.Namespace, order.Name)

			if err == nil {
				isIssuerRefPresent = true
			}
		}

		if order.Annotations != nil {
			if certName, found := order.Annotations[certmanagerv1.CertificateNameKey]; found && certName != "" {
				_, err = utils.GetCertificate(ctx, dynamicClient, order.Namespace, certName)
				if err == nil {
					isCertificatePresent = true
				}
			}
		}

		if !isIssuerRefPresent || !isCertificatePresent {
			staleOrders = append(staleOrders, order)
		}
	}

	return nil
}
