package utils

import (
	"github.com/cloudentity/acp-client-go/clients/hub/models"
	"github.com/go-json-experiment/json"
	"github.com/pkg/errors"
)

func FromTreeServerToPatch(data *models.TreeServer) (models.Rfc7396PatchOperation, error) {
	var (
		out = models.Rfc7396PatchOperation{}
		bts []byte
		err error
	)

	if bts, err = json.Marshal(data, json.FormatNilMapAsNull(true)); err != nil {
		return out, errors.Wrap(err, "failed to marshal tree server to yaml")
	}

	if err = json.Unmarshal(bts, &out); err != nil {
		return out, errors.Wrap(err, "failed to unmarshal yaml to patch")
	}

	return out, nil
}

func FromPatchToTreeServer(patch models.Rfc7396PatchOperation) (*models.TreeServer, error) {
	var (
		out = &models.TreeServer{}
		bts []byte
		err error
	)

	CleanPatch(patch)

	if bts, err = json.Marshal(patch, json.FormatNilMapAsNull(true)); err != nil {
		return out, errors.Wrap(err, "failed to marshal patch to yaml")
	}

	if err = json.Unmarshal(bts, out, json.RejectUnknownMembers(true)); err != nil {
		return out, errors.Wrap(err, "failed to unmarshal yaml to tree server")
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
		return out, errors.Wrap(err, "failed to marshal patch to yaml")
	}

	if err = json.Unmarshal(bts, &out); err != nil {
		return out, errors.Wrap(err, "failed to unmarshal yaml to patch")
	}

	return out, nil
}

// CleanPatch cleans fields that are available in system model but not available in hub model
func CleanPatch(patch models.Rfc7396PatchOperation) {
	delete(patch, "id")
	delete(patch, "tenant_id")
}

func FilterPatch(patch models.Rfc7396PatchOperation, filters []string) (models.Rfc7396PatchOperation, error) {
	var filterMap = map[string]bool{}

	if len(filters) == 0 {
		return patch, nil
	}

	for _, filter := range filters {
		filterMap[filter] = true
	}

	if !filterMap["policies"] {
		delete(patch, "policies")
	}

	if !filterMap["apis"] {
		delete(patch, "apis")
	}

	if !filterMap["scopes"] {
		delete(patch, "scopes_without_service")
	}

	if !filterMap["clients"] {
		delete(patch, "clients")
	}

	if !filterMap["webhooks"] {
		delete(patch, "webhooks")
	}

	if !filterMap["scripts"] {
		delete(patch, "scripts")
	}

	if !filterMap["services"] {
		delete(patch, "services")
	}

	if !filterMap["theme_binding"] {
		delete(patch, "theme_binding")
	}

	if !filterMap["servers_bindings"] {
		delete(patch, "servers_bindings")
	}

	if !filterMap["custom_apps"] {
		delete(patch, "custom_apps")
	}

	if !filterMap["pools"] {
		delete(patch, "pools")
	}

	if !filterMap["ciba"] {
		delete(patch, "ciba_authentication_service")
	}

	if !filterMap["policy_execution_points"] {
		delete(patch, "policy_execution_points")
	}

	if !filterMap["script_execution_points"] {
		delete(patch, "script_execution_points")
	}

	if !filterMap["gateways"] {
		delete(patch, "gateways")
	}

	if !filterMap["server_consent"] {
		delete(patch, "server_consent")
	}

	if !filterMap["idps"] {
		delete(patch, "idps")
	}

	if !filterMap["claims"] {
		delete(patch, "claims")
	}

	return patch, nil
}
