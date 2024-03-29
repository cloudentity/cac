package data

import (
	"github.com/cloudentity/acp-client-go/clients/hub/models"
	"github.com/cloudentity/cac/internal/cac/utils"
	"github.com/go-openapi/strfmt"
)

type ServerValidator struct{}

var _ ValidatorApi = &ServerValidator{}

func (sv *ServerValidator) Validate(data *models.Rfc7396PatchOperation) error {
	var (
		err  error
		serv *models.TreeServer
	)
	if serv, err = utils.FromPatchToModel[models.TreeServer](*data); err != nil {
		return err
	}

	if err = serv.Validate(strfmt.Default); err != nil {
		return err
	}

	return nil
}
