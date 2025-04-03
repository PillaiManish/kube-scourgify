package utils

import (
	apiv1 "github.com/openshift/api/operator/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StaleResourceSpec struct {
	Group             string   `json:"group"`
	Version           string   `json:"version"`
	Kind              string   `json:"kind"`
	StaleResourceList []string `json:"staleResourceList"`
}

type StaleResourceStatus struct {
	apiv1.OperatorStatus `json:",inline"`
}

type StaleResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StaleResourceSpec   `json:"spec,omitempty"`
	Status StaleResourceStatus `json:"status,omitempty"`
}

type StaleResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []StaleResource `json:"items"`
}
