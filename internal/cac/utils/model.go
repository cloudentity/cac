package utils

import (
	"github.com/cloudentity/acp-client-go/clients/hub/models"
	"github.com/go-json-experiment/json"
	"github.com/pkg/errors"
)

func FromModelToPatch[T any](data *T) (models.Rfc7396PatchOperation, error) {
	var (
		out = models.Rfc7396PatchOperation{}
		bts []byte
		err error
	)

	if bts, err = json.Marshal(data, json.FormatNilMapAsNull(true)); err != nil {
		return out, errors.Wrapf(err, "failed to marshal %T to yaml", out)
	}

	if err = json.Unmarshal(bts, &out); err != nil {
		return out, errors.Wrap(err, "failed to unmarshal yaml to patch")
	}

	return out, nil
}

func FromPatchToModel[T any](patch models.Rfc7396PatchOperation) (*T, error) {
	return FromPatchToModelWithOptions[T](patch, json.RejectUnknownMembers(true))
}

func FromPatchToModelWithOptions[T any](patch models.Rfc7396PatchOperation, unOpts... json.Options) (*T, error) {
	var (
		out = new(T)
		bts []byte
		err error
	)

	CleanPatch(patch)

	if bts, err = json.Marshal(patch, json.FormatNilMapAsNull(true)); err != nil {
		return out, errors.Wrap(err, "failed to marshal patch to json")
	}

	if err = json.Unmarshal(bts, out, unOpts...); err != nil {
		return out, errors.Wrapf(err, "failed to unmarshal json to %T", out)
	}

	return out, nil
}

func NormalizePatch(patch models.Rfc7396PatchOperation) (models.Rfc7396PatchOperation, error) {
	var (
		out = models.Rfc7396PatchOperation{}
		bts []byte
		err error
	)

	if bts, err = json.Marshal(patch, json.FormatNilMapAsNull(true)); err != nil {
		return out, errors.Wrap(err, "failed to marshal patch to json")
	}

	if err = json.Unmarshal(bts, &out); err != nil {
		return out, errors.Wrap(err, "failed to unmarshal json to patch")
	}

	return out, nil
}

func CleanPatch(patch models.Rfc7396PatchOperation) {
	// clean fields that are available in system model but not available in hub model
	delete(patch, "id")
	delete(patch, "tenant_id")
}

var staticFilterMappings = map[string]string{
	"scopes": "scopes_without_service",
	"ciba":   "ciba_authentication_service",
}

func FilterPatch(patch models.Rfc7396PatchOperation, filters []string) (models.Rfc7396PatchOperation, error) {
	if len(filters) == 0 {
		return patch, nil
	}

	var newPatch = models.Rfc7396PatchOperation{}

	for _, filter := range filters {
		if mapped, ok := staticFilterMappings[filter]; ok {
			filter = mapped
		}

		if _, ok := patch[filter]; ok {
			newPatch[filter] = patch[filter]
		}
	}

	return newPatch, nil
}
