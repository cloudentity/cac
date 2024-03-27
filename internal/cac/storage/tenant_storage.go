package storage

import (
	"context"
	"github.com/cloudentity/acp-client-go/clients/hub/models"
	smodels "github.com/cloudentity/acp-client-go/clients/system/models"
	"github.com/cloudentity/cac/internal/cac/api"
	"github.com/cloudentity/cac/internal/cac/utils"
	"github.com/go-json-experiment/json"
	"path/filepath"
)

func InitTenantStorage(config *Configuration) Storage {
	return &TenantStorage{
		Config:        config,
		ServerStorage: InitServerStorage(config),
	}
}

type TenantStorage struct {
	Config        *Configuration
	ServerStorage Storage
}

func (t *TenantStorage) Write(ctx context.Context, data models.Rfc7396PatchOperation, opts ...api.SourceOpt) error {
	var (
		path  = t.Config.DirPath
		model *models.TreeTenant
		err   error
	)

	if model, err = utils.FromPatchToModel[models.TreeTenant](data); err != nil {
		return err
	}

	if err = t.storeTenant(*model); err != nil {
		return err
	}

	if err = writeFiles(model.Pools,
		filepath.Join(path, "pools"),
		func(id string, it models.TreePool) string { return it.Name }); err != nil {
		return err
	}

	if err = writeFiles(model.Schemas,
		filepath.Join(path, "schemas"),
		func(id string, it models.TreeSchema) string { return it.Name }); err != nil {
		return err
	}

	if err = writeFiles(model.MfaMethods,
		filepath.Join(path, "mfa_methods"),
		func(id string, it models.TreeMFAMethod) string { return it.Mechanism }); err != nil {
		return err
	}

	if err = writeFiles(model.Themes,
		filepath.Join(path, "themes"),
		func(id string, it models.TreeTheme) string { return it.Name }); err != nil {
		return err
	}

	for k, server := range model.Servers {
		opts = append(opts, api.WithWorkspace(k))
		var serverData models.Rfc7396PatchOperation
		if serverData, err = utils.FromModelToPatch(&server); err != nil {
			return err
		}

		if err = t.ServerStorage.Write(ctx, serverData, opts...); err != nil {
			return err
		}
	}

	return nil
}

func (t *TenantStorage) Read(ctx context.Context, opts ...api.SourceOpt) (models.Rfc7396PatchOperation, error) {
	var (
		path       = t.Config.DirPath
		tenant     models.Rfc7396PatchOperation
		options    = &api.Options{}
		workspaces []string
		err        error
	)

	for _, opt := range opts {
		opt(options)
	}

	if tenant, err = readFile(filepath.Join(path, "tenant")); err != nil {
		return nil, err
	}

	if err = readFilesToMap(tenant, "pools", filepath.Join(path, "pools")); err != nil {
		return nil, err
	}

	if err = readFilesToMap(tenant, "schemas", filepath.Join(path, "schemas")); err != nil {
		return nil, err
	}

	if err = readFilesToMap(tenant, "mfa_methods", filepath.Join(path, "mfa_methods")); err != nil {
		return nil, err
	}

	if err = readFilesToMap(tenant, "themes", filepath.Join(path, "themes")); err != nil {
		return nil, err
	}

	if workspaces, err = listDirsInPath(filepath.Join(path, "workspaces")); err != nil {
		return nil, err
	}

	if len(workspaces) > 0 {
		var servers = map[string]any{}

		for _, workspace := range workspaces {
			var workspaceConfig models.Rfc7396PatchOperation
			opts = append(opts, api.WithWorkspace(workspace))

			if workspaceConfig, err = t.ServerStorage.Read(ctx, opts...); err != nil {
				return nil, err
			}

			id := workspaceConfig["id"].(string)
			delete(workspaceConfig, "id")
			delete(workspaceConfig, "tenant_id")
			servers[id] = workspaceConfig
		}

		tenant["servers"] = servers
	}

	if tenant, err = utils.FilterPatch(tenant, options.Filters); err != nil {
		return nil, err
	}

	return tenant, nil
}

// storeTenant stores the tenant data in the file
// it accepts a struct directly to make sure we only modify a copy
func (t *TenantStorage) storeTenant(data models.TreeTenant) error {
	var (
		path   = filepath.Join(t.Config.DirPath, "tenant")
		tenant = smodels.TenantDump{}
		bts    []byte
		err    error
	)

	data.MfaMethods = nil
	data.Themes = nil
	data.Servers = nil

	// serialize the tenant data into system/models to remove the dependencies which are stored in separate files
	if bts, err = json.Marshal(data); err != nil {
		return err
	}

	if err = json.Unmarshal(bts, &tenant); err != nil {
		return err
	}

	if err = writeFile(tenant, path); err != nil {
		return err
	}

	return nil
}

var _ Storage = &TenantStorage{}
