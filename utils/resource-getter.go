package utils

import (
	"context"
	"fmt"
	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

func GetSecretByName(ctx context.Context, kubeClient kubernetes.Interface, secretName string) (*corev1.Secret, error) {
	secret, err := kubeClient.CoreV1().Secrets("").Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if secret != nil {
		return secret, nil
	}

	return nil, fmt.Errorf("failed to fetch secret: %s", secretName)
}

func GetCertificate(ctx context.Context, dynamicClient dynamic.Interface, namespace, certificateName string) (*certmanagerv1.Certificate, error) {
	gvr := certmanagerv1.SchemeGroupVersion.WithResource("certificates")

	var certificate *certmanagerv1.Certificate

	unstructuredCertificate, err := dynamicClient.Resource(gvr).Namespace(namespace).Get(ctx, certificateName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredCertificate.Object, &certificate)
	if err != nil {
		return nil, err
	}

	return certificate, nil
}

func GetIssuerByName(ctx context.Context, dynamicClient dynamic.Interface, kind, namespace, issuerName string) (*certmanagerv1.Issuer, *certmanagerv1.ClusterIssuer, error) {
	gvr := schema.GroupVersionResource{}

	if kind == certmanagerv1.IssuerKind {
		gvr = certmanagerv1.SchemeGroupVersion.WithResource("issuers")
	} else if kind == certmanagerv1.ClusterIssuerKind {
		gvr = certmanagerv1.SchemeGroupVersion.WithResource("clusterissuers")
		namespace = ""
	}

	unStructuredIssuer, err := dynamicClient.Resource(gvr).Namespace(namespace).Get(ctx, issuerName, metav1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}

	switch kind {

	case certmanagerv1.IssuerKind:
		var issuer *certmanagerv1.Issuer

		err = runtime.DefaultUnstructuredConverter.FromUnstructured(unStructuredIssuer.UnstructuredContent(), issuer)
		if err != nil {
			return nil, nil, err
		}

		if unStructuredIssuer != nil {
			return issuer, nil, nil
		}

		return nil, nil, fmt.Errorf("failed to fetch issuer: %s", issuerName)
	case certmanagerv1.ClusterIssuerKind:
		clusterIssuer := &certmanagerv1.ClusterIssuer{}

		err = runtime.DefaultUnstructuredConverter.FromUnstructured(unStructuredIssuer.UnstructuredContent(), clusterIssuer)
		if err != nil {
			return nil, nil, err
		}

		if unStructuredIssuer != nil {
			return nil, clusterIssuer, nil
		}

		return nil, nil, fmt.Errorf("failed to fetch issuer: %s", issuerName)
	}

	return nil, nil, err

}

func GetCertificateRequestsByName(ctx context.Context, dynamicClient dynamic.Interface, namespace, certificateRequestName string) (*certmanagerv1.CertificateRequest, error) {
	gvr := certmanagerv1.SchemeGroupVersion.WithResource("certificaterequests")
	var certificateRequest *certmanagerv1.CertificateRequest

	unStructuredCertificateRequest, err := dynamicClient.Resource(gvr).Namespace(namespace).Get(ctx, certificateRequestName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unStructuredCertificateRequest.Object, &certificateRequest)
	if err != nil {
		return nil, err
	}

	return certificateRequest, nil
}
