package challenges

import (
	"context"
	acmev1 "github.com/cert-manager/cert-manager/pkg/apis/acme/v1"
	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"kube-scourgify/utils"
)

func FindStaleChallenges(ctx context.Context, dynamicClient *dynamic.DynamicClient, conditions utils.Conditions, deleteFlag bool) error {
	var staleChallenges []*acmev1.Challenge

	gvr := certmanagerv1.SchemeGroupVersion.WithResource(utils.CHALLENGES)

	unStructuredChallenges, err := dynamicClient.Resource(gvr).List(ctx, v1.ListOptions{})
	if err != nil {
		return err
	}

	for _, unStructuredChallenge := range unStructuredChallenges.Items {
		challenge := &acmev1.Challenge{}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(unStructuredChallenge.UnstructuredContent(), challenge)
		if err != nil {
			return err
		}

		if challenge.Spec.IssuerRef.Group == acmev1.SchemeGroupVersion.Group {
			_, _, err = utils.GetIssuerByName(ctx, dynamicClient, challenge.Spec.IssuerRef.Kind, challenge.Namespace, challenge.Spec.IssuerRef.Name)

			if err != nil {
				if errors.IsNotFound(err) {
					staleChallenges = append(staleChallenges, challenge)
				}
			}
		}

	}

	if deleteFlag {
		for _, challenge := range staleChallenges {
			err = dynamicClient.Resource(gvr).Namespace(challenge.Namespace).Delete(ctx, challenge.Name, v1.DeleteOptions{})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
