package data

import (
	"github.com/cloudentity/acp-client-go/clients/hub/models"
	"github.com/cloudentity/cac/internal/cac/utils"
	"github.com/go-openapi/strfmt"
)

type TenantValidator struct{}

var _ ValidatorApi = &TenantValidator{}

func (sv *TenantValidator) Validate(data *models.Rfc7396PatchOperation) error {
	var (
		err    error
		tenant *models.TreeTenant
	)
	if tenant, err = utils.FromPatchToModel[models.TreeTenant](*data); err != nil {
		return err
	}

	if err = tenant.Validate(strfmt.Default); err != nil {
		return err
	}

	return nil
}
