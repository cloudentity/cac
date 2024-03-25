package cac

import (
	"github.com/cloudentity/cac/internal/cac/api"
	"github.com/cloudentity/cac/internal/cac/client"
	"github.com/cloudentity/cac/internal/cac/config"
	"github.com/cloudentity/cac/internal/cac/logging"
	"github.com/cloudentity/cac/internal/cac/storage"
	"golang.org/x/exp/slog"
	"strings"
)

type Application struct {
	Config     *config.Configuration
	RootConfig *config.RootConfiguration
	Client     *client.Client
	Storage    storage.Storage
}

func InitApp(configPath string, profile string, tenant bool) (app *Application, err error) {
	app = &Application{}

	if app.RootConfig, err = config.InitConfig(configPath); err != nil {
		return app, err
	}

	if app.Config, err = app.RootConfig.ForProfile(profile); err != nil {
		return app, err
	}

	if err = logging.InitLogging(app.Config.Logging); err != nil {
		return app, err
	}

	slog.Debug("config", "c", app.Config.Client)

	if app.Config.Client != nil {
		if app.Client, err = client.InitClient(app.Config.Client); err != nil {
			return app, err
		}
	}

	var constructor = storage.InitServerStorage

	if tenant {
		constructor = storage.InitTenantStorage
	}

	if app.Config.Storage != nil {
		if app.Storage, err = storage.InitMultiStorage(app.Config.Storage, constructor); err != nil {
			return app, err
		}
	}

	slog.Info("Initiated application")

	return app, nil
}

func (a *Application) PickSource(source string, tenant bool) (api.Source, error) {
	var (
		conf       *config.Configuration
		sourceType api.SourceType
		err        error

		profile, sourceS, found = strings.Cut(source, "@")
	)

	if !found {
		sourceS = profile
		profile = a.Config.Name

		slog.With("profile", profile).With("source", sourceS).Debug("@ not found in source, using default profile")
	}

	if conf, err = a.RootConfig.ForProfile(profile); err != nil {
		return nil, err
	}

	if sourceType, err = api.SourceFromString(sourceS); err != nil {
		return nil, err
	}

	var constructor = storage.InitServerStorage

	if tenant {
		constructor = storage.InitTenantStorage
	}

	switch sourceType {
	case api.SourceLocal:
		return storage.InitMultiStorage(conf.Storage, constructor)
	case api.SourceRemote:
		if c, err := client.InitClient(conf.Client); err != nil {
			return nil, err
		} else {
			if tenant {
				return c.Tenant(), nil
			}

			return c, nil
		}
	}

	return nil, api.ErrUnknownSource
}
