package storage

import (
	"encoding/json"
	adminmodels "github.com/cloudentity/acp-client-go/clients/admin/models"
	"github.com/cloudentity/acp-client-go/clients/hub/models"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slog"
	"path/filepath"
)

type Configuration struct {
	DirPath string `json:"dir_path"`
}

var DefaultConfig = Configuration{
	DirPath: "data",
}

func InitStorage(config Configuration) *SingleStorage {
	return &SingleStorage{
		Config: config,
	}
}

type Storage interface {
	Store(workspace string, data *models.TreeServer) error
	Read(workspace string) (models.TreeServer, error)
}

type SingleStorage struct {
	Config Configuration
}

var _ Storage = &SingleStorage{}

func (s *SingleStorage) Store(workspace string, data *models.TreeServer) error {
	var (
		workspacePath = s.workspacePath(workspace)
		err           error
	)

	if err = s.storeServer(workspace, data); err != nil {
		return err
	}

	if err = writeFiles(data.Clients,
		filepath.Join(workspacePath, "clients"),
		func(id string, it models.TreeClient) string { return it.ClientName }); err != nil {
		return err
	}

	if err = writeFiles(data.Idps,
		filepath.Join(workspacePath, "idps"),
		func(id string, it models.TreeIDP) string { return it.Name }); err != nil {
		return err
	}

	if err = writeFile(data.Claims, filepath.Join(workspacePath, "claims")); err != nil {
		return err
	}

	if err = writeFiles(data.CustomApps,
		filepath.Join(workspacePath, "custom_apps"),
		func(id string, it models.TreeCustomApp) string { return it.Name }); err != nil {
		return err
	}

	if err = writeFiles(data.Gateways,
		filepath.Join(workspacePath, "gateways"),
		func(id string, it models.TreeGateway) string { return it.Name }); err != nil {
		return err
	}

	if err = writeFile(data.PolicyExecutionPoints, filepath.Join(workspacePath, "policy_execution_points")); err != nil {
		return err
	}

	if err = writeFiles(data.Pools,
		filepath.Join(workspacePath, "pools"),
		func(id string, it models.TreePool) string { return it.Name }); err != nil {
		return err
	}

	if err = writeFile(data.ScopesWithoutService, filepath.Join(workspacePath, "scopes")); err != nil {
		return err
	}

	if err = writeFile(data.ScriptExecutionPoints, filepath.Join(workspacePath, "script_execution_points")); err != nil {
		return err
	}

	if err = writeFile(data.ServerConsent, filepath.Join(workspacePath, "consent")); err != nil {
		return err
	}

	if len(data.ServersBindings) > 0 {
		if err = writeFile(map[string]any{
			"bindings": maps.Keys(data.ServersBindings),
		}, filepath.Join(workspacePath, "server_bindings")); err != nil {
			return err
		}
	}

	if err = writeFiles(data.Services,
		filepath.Join(workspacePath, "services"),
		func(id string, it models.TreeService) string { return it.Name }); err != nil {
		return err
	}

	if data.ThemeBinding != nil && data.ThemeBinding.ThemeID != "" {
		if err = writeFile(data.ThemeBinding, filepath.Join(workspacePath, "theme_binding")); err != nil {
			return err
		}
	}

	if err = writeFiles(data.Webhooks,
		filepath.Join(workspacePath, "webhooks"),
		func(id string, it models.TreeWebhook) string { return id }); err != nil {
		return err
	}

	if err = storeScripts(data.Scripts, filepath.Join(workspacePath, "scripts")); err != nil {
		return err
	}

	if err = StorePolicies(data.Policies, filepath.Join(workspacePath, "policies")); err != nil {
		return err
	}

	slog.Info("Workspace configuration successfully stored", "workspace", workspace, "path", workspacePath)

	return nil
}

func (s *SingleStorage) Read(workspace string) (models.TreeServer, error) {
	var (
		server = models.TreeServer{
			Clients:              models.TreeClients{},
			Idps:                 models.TreeIDPs{},
			CustomApps:           models.TreeCustomApps{},
			Gateways:             models.TreeGateways{},
			Pools:                models.TreePools{},
			Webhooks:             models.TreeWebhooks{},
			Scripts:              models.TreeScripts{},
			Policies:             models.TreePolicies{},
			Services:             models.TreeServices{},
			ScopesWithoutService: models.TreeScopes{},
			ServersBindings:      models.TreeServersBindings{},
		}
		path = s.workspacePath(workspace)
		err  error
	)

	if err = readFile(filepath.Join(path, "server"), &server); err != nil {
		return server, err
	}

	if err = readFiles(filepath.Join(path, "clients"), &server.Clients); err != nil {
		return server, err
	}

	if err = readFiles(filepath.Join(path, "idps"), &server.Idps); err != nil {
		return server, err
	}

	if err = readFile(filepath.Join(path, "claims"), &server.Claims); err != nil {
		return server, err
	}

	if err = readFiles(filepath.Join(path, "custom_apps"), &server.CustomApps); err != nil {
		return server, err
	}

	if err = readFiles(filepath.Join(path, "gateways"), &server.Gateways); err != nil {
		return server, err
	}

	if err = readFile(filepath.Join(path, "policy_execution_points"), &server.PolicyExecutionPoints); err != nil {
		return server, err
	}

	if err = readFiles(filepath.Join(path, "pools"), &server.Pools); err != nil {
		return server, err
	}

	if err = readFile(filepath.Join(path, "scopes"), &server.ScopesWithoutService); err != nil {
		return server, err
	}

	if err = readFile(filepath.Join(path, "script_execution_points"), &server.ScriptExecutionPoints); err != nil {
		return server, err
	}

	if err = readFile(filepath.Join(path, "consent"), server.ServerConsent, func(opts *ReadFileOpts[models.TreeServerConsent]) {
		opts.Constructor = func() *models.TreeServerConsent {
			server.ServerConsent = &models.TreeServerConsent{}
			return server.ServerConsent
		}
	}); err != nil {
		return server, err
	}

	sb := map[string]any{}
	if err = readFile(filepath.Join(path, "server_bindings"), &sb); err != nil {
		return server, err
	}

	if bindings, ok := sb["bindings"].([]any); ok && len(bindings) != 0 {
		server.ServersBindings = models.TreeServersBindings{}

		for _, binding := range bindings {
			server.ServersBindings[binding.(string)] = true
		}
	}

	if err = readFiles(filepath.Join(path, "services"), &server.Services); err != nil {
		return server, err
	}

	if err = readFile(filepath.Join(path, "theme_binding"), server.ThemeBinding, func(opts *ReadFileOpts[models.TreeThemeBinding]) {
		opts.Constructor = func() *models.TreeThemeBinding {
			server.ThemeBinding = &models.TreeThemeBinding{}
			return server.ThemeBinding
		}
	}); err != nil {
		return server, err
	}

	if err = readFiles(filepath.Join(path, "webhooks"), &server.Webhooks); err != nil {
		return server, err
	}

	if err = readFiles(filepath.Join(path, "scripts"), &server.Scripts); err != nil {
		return server, err
	}

	if err = readFiles(filepath.Join(path, "policies"), &server.Policies); err != nil {
		return server, err
	}

	return server, nil
}

func (s *SingleStorage) workspacePath(workspace string) string {
	return filepath.Join(s.Config.DirPath, "workspaces", workspace)
}

func (s *SingleStorage) storeServer(workspace string, data *models.TreeServer) error {
	var (
		path   = filepath.Join(s.workspacePath(workspace), "server")
		server adminmodels.Server
		bts    []byte
		err    error
	)

	// serialize the server data into adminmodel to remove the dependencies which are stored in separate files
	if bts, err = json.Marshal(data); err != nil {
		return err
	}

	if err = json.Unmarshal(bts, &server); err != nil {
		return err
	}

	if err = writeFile(NewWithID(workspace, server), path); err != nil {
		return err
	}

	return nil
}
