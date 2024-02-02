package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/cloudentity/acp-client-go"
	"github.com/cloudentity/acp-client-go/clients/hub/client/workspace_configuration"
	"github.com/cloudentity/acp-client-go/clients/hub/models"
	"github.com/cloudentity/cac/internal/cac/api"
	"github.com/go-openapi/runtime"
	"golang.org/x/exp/slog"
	"net/http"
	"strings"
)

type Client struct {
	acp *acpclient.Client
}

func WithSecrets(secrets bool) api.SourceOpt {
	return func(o *api.Options) {
		o.Secrets = secrets
	}
}

var _ api.Source = &Client{}

func InitClient(config *Configuration) (c *Client, err error) {
	var (
		acp acpclient.Client
	)

	if config.Insecure {
		config.HttpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	}

	if acp, err = acpclient.New(config.Config); err != nil {
		return nil, err
	}

	slog.With("client", acp).Debug("Initiated client")

	return &Client{
		acp: &acp,
	}, nil
}

func (c *Client) Read(ctx context.Context, workspace string, opts ...api.SourceOpt) (*models.TreeServer, error) {
	var (
		ok      *workspace_configuration.ExportWorkspaceConfigOK
		options = &api.Options{}
		err     error
	)

	for _, opt := range opts {
		opt(options)
	}

	if ok, err = c.acp.Hub.WorkspaceConfiguration.
		ExportWorkspaceConfig(workspace_configuration.
			NewExportWorkspaceConfigParams().
			WithContext(ctx).
			WithWithCredentials(&options.Secrets).
			WithWid(workspace), nil); err != nil {
		return nil, err
	}

	return ok.Payload, nil
}

func (c *Client) Write(ctx context.Context, workspace string, input *models.TreeServer) error {
	var (
		mode = "update"
		err  error
	)

	slog.Info("writing workspace configuration", "in", input.Name)

	if _, err = c.acp.Hub.WorkspaceConfiguration.
		ImportWorkspaceConfig(workspace_configuration.
			NewImportWorkspaceConfigParams().
			WithContext(ctx).
			WithWid(workspace).
			WithMode(&mode).
			WithConfig(input), nil, func(operation *runtime.ClientOperation) {
			operation.PathPattern = "/workspaces/{wid}/promote/config"
		}); err != nil {

		// temporary workaround for acp not returning the advertised 204 no content status code
		if strings.Contains(err.Error(), "status 200") {
			slog.Info(err.Error())
			return nil
		}

		return err
	}

	return nil
}

func (c *Client) String() string {
	return fmt.Sprintf("client: %v", c.acp.Config.IssuerURL)
}
