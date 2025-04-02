package utils

import "os"

const (
	//TODO: RESOURCE_NAME_KEY      = "resourceName"
	RESOURCE_VERSION_KEY   = "resourceVersion"
	RESOURCE_GROUP_KEY     = "resourceGroup"
	RESOURCE_KIND_KEY      = "resourceKind"
	RESOURCE_NAMESPACE_KEY = "resourceNamespace"
	CONDITIONS_FILEPATH    = "conditions"
	SCOUR_VERSION          = "1.0.0"
	CERTIFICATE            = "certificate"
	CERTIFICATES           = "certificates"
	ISSUER                 = "issuer"
	ISSUERS                = "issuers"
	CLUSTER_ISSUER         = "clusterissuer"
	CLUSTER_ISSUERS        = "clusterissuers"
	RESOURCE_CR_NAME       = "resourceCRName"
)

var (
	DEFAULT_CONFIG_PATH = os.Getenv("HOME") + "/.kube/config"
)
