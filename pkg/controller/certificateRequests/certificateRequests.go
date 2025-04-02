package certificateRequests

import (
	"kube-scourgify/utils"
	"strings"
)

func FindStaleCertificateRequests(conditions utils.Conditions) error {

	for i := range conditions.Deps {
		if strings.ToLower(conditions.Deps[i]) == utils.CERTIFICATE || strings.ToLower(conditions.Deps[i]) == utils.CERTIFICATES {
			//Call corresponding function
		} else if strings.ToLower(conditions.Deps[i]) == utils.ISSUER || strings.ToLower(conditions.Deps[i]) == utils.ISSUERS {
			//Call corresponding function
		}
	}
	return nil
}
