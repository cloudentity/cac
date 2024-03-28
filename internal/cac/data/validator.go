package data

import "github.com/cloudentity/acp-client-go/clients/hub/models"

type ValidatorApi interface {
	Validate(data *models.Rfc7396PatchOperation) error
}

func CreateValidator(tenant bool) ValidatorApi {
	if tenant {
		return &TenantValidator{}
	}

	return &ServerValidator{}
}
