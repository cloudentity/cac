package data

import (
	"github.com/cloudentity/acp-client-go/clients/hub/models"
	"github.com/cloudentity/cac/internal/cac/api"
	"github.com/cloudentity/cac/internal/cac/utils"
	"github.com/go-openapi/strfmt"
)

type ServerValidator struct{}

var _ ValidatorApi = &ServerValidator{}

func (sv *ServerValidator) Validate(data api.PatchInterface) error {
	var (
		err   error
		serv  *models.TreeServer
		sdata = data.GetData()
	)

	utils.CleanPatch(sdata)

	if serv, err = utils.FromPatchToModel[models.TreeServer](sdata); err != nil {
		return err
	}

	allowToDeleteScriptExecutionPoints(serv)

	if err = serv.Validate(strfmt.Default); err != nil {
		return err
	}

	return nil
}
