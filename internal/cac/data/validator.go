package data

import (
	"github.com/cloudentity/cac/internal/cac/api"
)

type ValidatorApi interface {
	Validate(data api.Patch) error
}
