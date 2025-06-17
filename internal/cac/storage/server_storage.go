package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/cloudentity/acp-client-go/clients/hub/models"
	smodels "github.com/cloudentity/acp-client-go/clients/system/models"
	"github.com/cloudentity/cac/internal/cac/api"
	"github.com/cloudentity/cac/internal/cac/utils"
	"github.com/pkg/errors"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slog"
)

type Configuration struct {
	DirPath string `json:"dir_path"`
}

var DefaultConfig = Configuration{
	DirPath: "data",
}

func InitServerStorage(config *Configuration) Storage {
	return &ServerStorage{
		Config: config,
	}
}

type ServerStorage struct {
	Config *Configuration
}

var _ Storage = &ServerStorage{}
var _ api.Source = &ServerStorage{}

func (s *ServerStorage) Write(ctx context.Context, input api.Patch, opts ...api.SourceOpt) error {
	var (
		workspacePath string
		workspace     string
		data          *models.TreeServer
		options       = &api.Options{}
		err           error
	)

	for _, opt := range opts {
		opt(options)
	}

	slog.Debug("write server data", "options", options)

	if workspace = options.Workspace; workspace == "" {
		return errors.New("workspace is required to write to server storage")
	}

	workspacePath = s.workspacePath(workspace)

	if data, err = utils.FromPatchToModel[models.TreeServer](input.GetData()); err != nil {
		return errors.Wrap(err, "failed to convert patch to tree server")
	}

	if err = s.storeServer(workspace, data); err != nil {
		return errors.Wrapf(err, "failed to store server data for workspace %s", workspace)
	}

	if err = writeFiles(data.Clients,
		filepath.Join(workspacePath, "clients"),
		func(id string, it models.TreeClient) string { return it.ClientName }); err != nil {
		return errors.Wrapf(err, "failed to write clients for workspace %s", workspace)
	}

	if err = writeFiles(data.Idps,
		filepath.Join(workspacePath, "idps"),
		func(id string, it models.TreeIDP) string { return it.Name }); err != nil {
		return errors.Wrapf(err, "failed to write idps for workspace %s", workspace)
	}

	if err = writeFile(data.Claims, filepath.Join(workspacePath, "claims")); err != nil {
		return errors.Wrapf(err, "failed to write claims for workspace %s", workspace)
	}

	if err = writeFiles(data.CustomApps,
		filepath.Join(workspacePath, "custom_apps"),
		func(id string, it models.TreeCustomApp) string { return it.Name }); err != nil {
		return errors.Wrapf(err, "failed to write custom apps for workspace %s", workspace)
	}

	if err = writeFiles(data.Gateways,
		filepath.Join(workspacePath, "gateways"),
		func(id string, it models.TreeGateway) string { return it.Name }); err != nil {
		return errors.Wrapf(err, "failed to write gateways for workspace %s", workspace)
	}

	if err = writeFile(data.PolicyExecutionPoints, filepath.Join(workspacePath, "policy_execution_points")); err != nil {
		return errors.Wrapf(err, "failed to write policy execution points for workspace %s", workspace)
	}

	if err = writeFiles(data.Pools,
		filepath.Join(workspacePath, "pools"),
		func(id string, it models.TreePool) string { return it.Name }); err != nil {
		return errors.Wrapf(err, "failed to write pools for workspace %s", workspace)
	}

	if err = writeFile(data.ScopesWithoutService, filepath.Join(workspacePath, "scopes")); err != nil {
		return errors.Wrapf(err, "failed to write scopes for workspace %s", workspace)
	}

	if err = writeFile(data.ScriptExecutionPoints, filepath.Join(workspacePath, "script_execution_points")); err != nil {
		return errors.Wrapf(err, "failed to write script execution points for workspace %s", workspace)
	}

	if err = writeFile(data.ServerConsent, filepath.Join(workspacePath, "consent")); err != nil {
		return errors.Wrapf(err, "failed to write server consent for workspace %s", workspace)
	}

	if len(data.ServersBindings) > 0 {
		if err = writeFile(map[string]any{
			"bindings": maps.Keys(data.ServersBindings),
		}, filepath.Join(workspacePath, "servers_bindings")); err != nil {
			return errors.Wrapf(err, "failed to write server bindings for workspace %s", workspace)
		}
	}

	if err = writeFiles(data.Services,
		filepath.Join(workspacePath, "services"),
		func(id string, it models.TreeService) string { return it.Name }); err != nil {
		return errors.Wrapf(err, "failed to write services for workspace %s", workspace)
	}

	if data.ThemeBinding != nil && data.ThemeBinding.ThemeID != "" {
		if err = writeFile(data.ThemeBinding, filepath.Join(workspacePath, "theme_binding")); err != nil {
			return errors.Wrapf(err, "failed to write theme binding for workspace %s", workspace)
		}
	}

	if err = writeFiles(data.Webhooks,
		filepath.Join(workspacePath, "webhooks"),
		func(id string, it models.TreeWebhook) string { return id }); err != nil {
		return errors.Wrapf(err, "failed to write webhooks for workspace %s", workspace)
	}

	if err = writeFile(data.CibaAuthenticationService, filepath.Join(workspacePath, "ciba")); err != nil {
		return errors.Wrapf(err, "failed to write ciba authentication service for workspace %s", workspace)
	}

	if err = storeScripts(data.Scripts, filepath.Join(workspacePath, "scripts")); err != nil {
		return errors.Wrapf(err, "failed to store scripts for workspace %s", workspace)
	}

	if err = StorePolicies(data.Policies, filepath.Join(workspacePath, "policies")); err != nil {
		return errors.Wrapf(err, "failed to store policies for workspace %s", workspace)
	}

	if options.Secrets {
		slog.Debug("trying to write secrets", "server", options.Workspace)
	}

	if options.Secrets {
		ext, ok := input.GetExtensions().(*api.ServerExtensions)

		if !ok {
			return errors.New("extensions are required to write secrets")
		}

		for _, secret := range ext.Secrets {
			secret.Secret = "" // clear the secret to avoid storing encrypted secrets in the storage
		}

		if err = writeFiles(ext.Secrets,
			filepath.Join(workspacePath, "secrets"),
			func(id string, it *smodels.Secret) string { return id }); err != nil {
			return errors.Wrapf(err, "failed to write secrets for workspace %s", workspace)
		}
	}

	slog.Info("Workspace configuration successfully stored", "workspace", workspace, "path", workspacePath)

	return nil
}

func (s *ServerStorage) Read(ctx context.Context, opts ...api.SourceOpt) (api.Patch, error) {
	var (
		path      string
		workspace string
		server    models.Rfc7396PatchOperation
		ext       = models.Rfc7396PatchOperation{}
		options   = &api.Options{}
		err       error
	)

	for _, opt := range opts {
		opt(options)
	}

	if workspace = options.Workspace; workspace == "" {
		return nil, errors.New("workspace is required to read from server storage")
	}

	path = s.workspacePath(workspace)

	if server, err = readFile(filepath.Join(path, "server")); err != nil {
		return nil, err
	}

	if err = readFilesToMap(server, "clients", filepath.Join(path, "clients")); err != nil {
		return nil, err
	}

	if err = readFilesToMap(server, "idps", filepath.Join(path, "idps")); err != nil {
		return nil, err
	}

	if err = readFileToMap(server, "claims", filepath.Join(path, "claims")); err != nil {
		return nil, err
	}

	if err = readFilesToMap(server, "custom_apps", filepath.Join(path, "custom_apps")); err != nil {
		return nil, err
	}

	if err = readFilesToMap(server, "gateways", filepath.Join(path, "gateways")); err != nil {
		return nil, err
	}

	if err = readFileToMap(server, "policy_execution_points", filepath.Join(path, "policy_execution_points")); err != nil {
		return nil, err
	}

	if err = readFilesToMap(server, "pools", filepath.Join(path, "pools")); err != nil {
		return nil, err
	}

	if err = readFileToMap(server, "scopes_without_service", filepath.Join(path, "scopes")); err != nil {
		return nil, err
	}

	if err = readFileToMap(server, "script_execution_points", filepath.Join(path, "script_execution_points")); err != nil {
		return nil, err
	}

	if err = readFileToMap(server, "server_consent", filepath.Join(path, "consent")); err != nil {
		return nil, err
	}

	if err = readFileToMap(server, "ciba_authentication_service", filepath.Join(path, "ciba")); err != nil {
		return nil, err
	}

	var sb map[string]any
	if sb, err = readFile(filepath.Join(path, "servers_bindings")); err != nil {
		return nil, err
	}

	if bindings, ok := sb["bindings"].([]any); ok && len(bindings) != 0 {
		binds := map[string]any{}

		for _, binding := range bindings {
			binds[binding.(string)] = true
		}

		server["servers_bindings"] = binds
	}

	if err = readFilesToMap(server, "services", filepath.Join(path, "services")); err != nil {
		return nil, err
	}

	if err = readFileToMap(server, "theme_binding", filepath.Join(path, "theme_binding")); err != nil {
		return nil, err
	}

	if err = readFilesToMap(server, "webhooks", filepath.Join(path, "webhooks")); err != nil {
		return nil, err
	}

	if err = readFilesToMap(server, "scripts", filepath.Join(path, "scripts")); err != nil {
		return nil, err
	}

	if err = readFilesToMap(server, "policies", filepath.Join(path, "policies")); err != nil {
		return nil, err
	}

	if err = readFilesToMap(ext, "secrets", filepath.Join(path, "secrets")); err != nil {
		return nil, err
	}

	if server, err = utils.FilterPatch(server, options.Filters); err != nil {
		return nil, err
	}

	sext, err := utils.FromPatchToModel[api.ServerExtensions](ext)

	if err != nil {
		return nil, errors.Wrap(err, "failed to convert extensions to model")
	}

	return &api.ServerPatch{
		Data: server,
		Ext:  sext,
	}, nil
}

func (s *ServerStorage) String() string {
	return fmt.Sprintf("server storage: %v", s.Config.DirPath)
}

func (s *ServerStorage) workspacePath(workspace string) string {
	return filepath.Join(s.Config.DirPath, "workspaces", workspace)
}

func (s *ServerStorage) storeServer(workspace string, data *models.TreeServer) error {
	var (
		path   = filepath.Join(s.workspacePath(workspace), "server")
		server smodels.ServerDump
		bts    []byte
		err    error
	)

	// serialize the server data into system/models to remove the dependencies which are stored in separate files
	if bts, err = json.Marshal(data); err != nil {
		return errors.Wrapf(err, "failed to marshal server data for workspace %s", workspace)
	}

	if err = json.Unmarshal(bts, &server); err != nil {
		return errors.Wrapf(err, "failed to unmarshal server data for workspace %s into system model", workspace)
	}

	server.ID = workspace

	if err = writeFile(server, path); err != nil {
		return err
	}

	return nil
}
