package data

import (
	"github.com/cloudentity/acp-client-go/clients/hub/models"
	"github.com/cloudentity/cac/internal/cac/api"
	"github.com/cloudentity/cac/internal/cac/utils"
	"github.com/go-openapi/strfmt"
)

type TenantValidator struct{}

var _ ValidatorApi = &TenantValidator{}

func (sv *TenantValidator) Validate(data api.PatchInterface) error {
	var (
		err    error
		tenant *models.TreeTenant
		tdata  = data.GetData()
	)

	utils.CleanPatch(tdata)

	if tenant, err = utils.FromPatchToModel[models.TreeTenant](tdata); err != nil {
		return err
	}

	for _, server := range tenant.Servers {
		allowToDeleteScriptExecutionPoints(&server)
	}

	if err = tenant.Validate(strfmt.Default); err != nil {
		return err
	}

	return nil
}
