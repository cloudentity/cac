package data

import "github.com/cloudentity/acp-client-go/clients/hub/models"

type ValidatorApi interface {
	Validate(data *models.Rfc7396PatchOperation) error
}
