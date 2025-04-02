package certificates

import (
	"kube-scourgify/utils"
	"strings"
)

func FindStaleCertificates(conditions utils.Conditions) error {

	for i := range conditions.Deps {
		if strings.ToLower(conditions.Deps[i]) == utils.ISSUER || strings.ToLower(conditions.Deps[i]) == utils.ISSUERS {
			//Call corresponding function
		} else if strings.ToLower(conditions.Deps[i]) == utils.CLUSTER_ISSUER || strings.ToLower(conditions.Deps[i]) == utils.CLUSTER_ISSUERS {
			//Call corresponding function
		}
	}
	return nil
}
