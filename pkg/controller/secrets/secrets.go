package secrets

import (
	"kube-scourgify/utils"
	"strings"
)

func FindStaleSecrets(conditions utils.Conditions) error {

	for i := range conditions.Deps {
		if strings.ToLower(conditions.Deps[i]) == utils.CERTIFICATE || strings.ToLower(conditions.Deps[i]) == utils.CERTIFICATES {
			//Call corresponding function
		}
	}
	return nil
}
